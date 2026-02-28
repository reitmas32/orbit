package chat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/fundathyon/base/core/customctx"
	"github.com/fundathyon/base/core/eventbus"
	"github.com/jalocompany/orbit/internal/core/auth"
	"github.com/jalocompany/orbit/internal/core/entities"
	"github.com/jalocompany/orbit/internal/core/ports"
)

const (
	MessageTypeDirect = "direct-message"
	MessageTypeGroup  = "group"
)

type Handler struct {
	service     ports.ChatService
	notifier    ports.UserNotifierService
	storage     ports.MessageStoreService    // Servicio para persistencia de mensajes
	multiplexer ports.MultiplexerService     // Servicio para multiplexado de grupos
	eb          eventbus.EventBus            // Bus de eventos para consumir desde RabbitMQ
	Mu          sync.Mutex                   // Mutex para sincronizar el acceso a las sesiones
	Sessions    map[string]*websocket.Conn   // Almacén de conexiones activas por SessionID
	Broadcast   chan entities.SessionMessage // Canal para transmitir mensajes a los clientes

	// Mapa de manejadores según el tipo de mensaje (direct-message, group, etc.)
	handlers map[string]func(ctx *customctx.CustomContext, msg entities.SessionMessage) error
}

func NewHandler(
	service ports.ChatService,
	notifier ports.UserNotifierService,
	storage ports.MessageStoreService,
	multiplexer ports.MultiplexerService,
	eb eventbus.EventBus,
) *Handler {
	h := &Handler{
		service:     service,
		notifier:    notifier,
		storage:     storage,
		multiplexer: multiplexer,
		eb:          eb,
		Sessions:    make(map[string]*websocket.Conn),
		Broadcast:   make(chan entities.SessionMessage, 100),
	}

	// 1. Inicializar mapa de manejadores de ruteo
	h.handlers = map[string]func(ctx *customctx.CustomContext, msg entities.SessionMessage) error{
		MessageTypeDirect: func(ctx *customctx.CustomContext, msg entities.SessionMessage) error {
			ctx.Logger().Info("[Direct] Publicando mensaje directo hacia el bus para %s", msg.RecipientID)
			return h.service.PublishMessage(ctx, msg)
		},
		MessageTypeGroup: func(ctx *customctx.CustomContext, msg entities.SessionMessage) error {
			ctx.Logger().Info("[Multiplexer] Redireccionando mensaje de grupo hacia el multiplexer")
			return h.multiplexer.Multiplex(ctx, msg)
		},
	}

	// Iniciar el consumidor de broadcast que enviará mensajes desde el canal interno hacia los WS
	go h.listenBroadcast()

	// Iniciar la suscripción a RabbitMQ que alimentará el canal Broadcast
	go h.startRabbitConsumer()

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow connections from any origin
	})
	if err != nil {
		log.Printf("failed to accept chat websocket connection: %v", err)
		return
	}

	ctx := h.enrichContext(r)

	// Recuperar los JWTClaims empaquetados desde el Auth Middleware
	var userID string
	if claims, ok := ctx.Value("claims").(*auth.JWTClaims); ok {
		ctx.Logger().Info("Nueva conexión de sesión para el usuario %s (AppID: %s)", claims.UserID, claims.AppID)
		userID = claims.UserID
	} else {
		ctx.Logger().Warn("No se encontraron los JWT claims dentro del context")
	}

	h.onConnect(ctx, c, userID)
	defer h.onDisconnect(ctx, c, userID)

	h.handleMessages(ctx, c, userID, r)
}

func (h *Handler) enrichContext(r *http.Request) *customctx.CustomContext {
	ctx, ok := r.Context().(*customctx.CustomContext)
	if !ok {
		log.Println("[Middleware: ChatHandler] Not enriched context")
		ctx = customctx.New(r.Context())
	}
	return ctx
}

func (h *Handler) handleMessages(ctx *customctx.CustomContext, c *websocket.Conn, userID string, r *http.Request) {
	for {
		var received entities.MessageReceived
		err := wsjson.Read(ctx.Context(), c, &received)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				return
			}
			ctx.Logger().Error("failed to read chat message: %v", err)
			return
		}

		received.CreatedAt = time.Now()
		received.SenderID = userID // El que envía el mensaje
		received.Text = fmt.Sprintf("[Recibido y enviado] %v", received.Text)

		// Construir el mensaje de sesión enriquecido para ser enviado al bus (RabbitMQ)
		// La SessionID en este caso es el TargetID para que el bus lo entregue al destinatario
		sessionMsg := entities.SessionMessage{
			RecipientID: received.TargetID, // DESTINATARIO: A quién queremos que le llegue
			Browser:     r.UserAgent(),
			IP:          r.RemoteAddr,
			SO:          "Unknown", // Metadata adicional que se puede extraer si es necesario
			Message:     received,
			TraceID:     fmt.Sprintf("%v", ctx.Value("trace_id")),
		}

		// 1. Persistir el mensaje primero (Core API o Mock) siempre
		if err := h.storage.SaveMessage(ctx, sessionMsg); err != nil {
			ctx.Logger().Error("Error al persistir mensaje: %v", err)
		}

		// 2. Ejecutar el manejador correspondiente según el tipo (direct-message, group, etc.)
		handler, exists := h.handlers[received.Type]
		if !exists {
			// Manejador por defecto (direct-message)
			handler = h.handlers[MessageTypeDirect]
		}

		if err := handler(ctx, sessionMsg); err != nil {
			ctx.Logger().Error("Error al procesar mensaje de tipo %s: %v", received.Type, err)
		}
	}
}

// startRabbitConsumer suscribe al handler al bus de eventos de RabbitMQ
func (h *Handler) startRabbitConsumer() {
	// Crear un contexto base para el consumidor
	ctx := customctx.New(context.Background())

	// El evento que queremos consumir (usamos una instancia vacía para el tipo)
	event := entities.SessionMessage{}

	// Suscribirse al bus
	deliveryRes := h.eb.Consume(ctx, event)
	if !deliveryRes.IsOk() {
		log.Printf("[Chat Handler] Error al iniciar consumidor de RabbitMQ: %v", deliveryRes.Error())
		return
	}

	deliveryChan := deliveryRes.Value()

	log.Println("[Chat Handler] Suscrito exitosamente a RabbitMQ para mensajes de broadcast")

	for msg := range deliveryChan {
		// Convertir el payload del evento a SessionMessage
		if sessionMsg, ok := msg.Event().Payload().(entities.SessionMessage); ok {
			// Enviar al canal interno de broadcast
			h.Broadcast <- sessionMsg
		} else {
			log.Printf("[Chat Handler] Payload de evento inválido recibido de RabbitMQ")
		}

		// Aceptar el mensaje en el bus
		_ = msg.Ack()
	}
}

// listenBroadcast escucha el canal interno para enviar mensajes a los WS correspondientes
func (h *Handler) listenBroadcast() {
	for msg := range h.Broadcast {
		h.Mu.Lock()
		conn, exists := h.Sessions[msg.RecipientID]
		h.Mu.Unlock()

		if exists {
			// Enviar el mensaje al cliente específico
			err := wsjson.Write(context.Background(), conn, msg.Message)
			if err != nil {
				log.Printf("failed to broadcast message to recipient %s: %v", msg.RecipientID, err)
			}
		}
	}
}

func (h *Handler) onConnect(ctx *customctx.CustomContext, c *websocket.Conn, userID string) {
	if userID == "" {
		return
	}

	h.Mu.Lock()
	h.Sessions[userID] = c
	h.Mu.Unlock()

	ctx.Logger().Info("Usuario %s registrado en la tabla de sesiones", userID)
	_ = h.notifier.NotifyConnect(ctx, userID)
}

func (h *Handler) onDisconnect(ctx *customctx.CustomContext, c *websocket.Conn, userID string) {
	if userID != "" {
		h.Mu.Lock()
		delete(h.Sessions, userID)
		h.Mu.Unlock()

		ctx.Logger().Info("Usuario %s eliminado de la tabla de sesiones", userID)
		_ = h.notifier.NotifyDisconnect(ctx, userID)
	}
	c.CloseNow()
}

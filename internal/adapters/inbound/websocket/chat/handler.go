package chat

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/auth"
	"github.com/jalocompany/orbit/internal/core/ports"
)

type Handler struct {
	service ports.ChatService
}

func NewHandler(service ports.ChatService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow connections from any origin
	})
	if err != nil {
		log.Printf("failed to accept chat websocket connection: %v", err)
		return
	}
	defer c.CloseNow()

	ctx := h.enrichContext(r)

	h.handleMessages(ctx, c)
}

func (h *Handler) enrichContext(r *http.Request) *customctx.CustomContext {
	ctx, ok := r.Context().(*customctx.CustomContext)

	// Intentar recuperar el context enriquecido e imprimir el log con Loguro
	if ok {
		ctx.Logger().Info("[Middleware: ChatHandler] Bypass ok")
	} else {
		log.Println("[Middleware: ChatHandler] Not enriched context")
		ctx = customctx.New(r.Context())
	}

	return ctx
}

func (h *Handler) handleMessages(ctx *customctx.CustomContext, c *websocket.Conn) {
	// Recuperar los JWTClaims empaquetados desde el Auth Middleware
	if claims, ok := ctx.Value("claims").(*auth.JWTClaims); ok {
		ctx.Logger().Info("Empezando a escuchar al usuario %s (AppID: %s)", claims.UserID, claims.AppID)
	} else {
		ctx.Logger().Warn("No se encontraron los JWT claims dentro del context")
	}

	for {
		var msg any

		// Read message unmarshaled from JSON
		err := wsjson.Read(ctx.Context(), c, &msg)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				return
			}
			ctx.Logger().Error("failed to read chat message: %v", err)
			return
		}

		ctx.Logger().Info("received chat message: %v", msg)

		// Use the port to process the message
		response := h.service.ProcessMessage(msg)

		// Write message back as JSON
		err = wsjson.Write(ctx.Context(), c, response)
		if err != nil {
			ctx.Logger().Error("failed to write chat message: %v", err)
			return
		}
	}
}

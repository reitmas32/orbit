package echo

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/jalocompany/orbit/internal/core/ports"
)

type Handler struct {
	service ports.EchoService
}

func NewHandler(service ports.EchoService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow connections from any origin
	})
	if err != nil {
		log.Printf("failed to accept websocket connection: %v", err)
		return
	}
	defer c.CloseNow()

	ctx := r.Context()

	for {
		var msg any

		// Read message unmarshaled from JSON
		err = wsjson.Read(ctx, c, &msg)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				return
			}
			log.Printf("failed to read message: %v", err)
			return
		}

		// Use the port to process the message (Hexagonal Architecture)
		response := h.service.ProcessMessage(msg)

		// Write message back as JSON
		err = wsjson.Write(ctx, c, response)
		if err != nil {
			log.Printf("failed to write message: %v", err)
			return
		}
	}
}

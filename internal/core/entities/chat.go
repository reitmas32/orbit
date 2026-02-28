package entities

import (
	"time"

	"github.com/fundathyon/base/core/eventbus"
)

// MessageReceived representa el contenido del mensaje enviado por el usuario
type MessageReceived struct {
	Text      string    `json:"text" bson:"text"`
	Type      string    `json:"type" bson:"type"` // e.g., "chat", "image", "system"
	SenderID  string    `json:"sender_id" bson:"sender_id"`
	TargetID  string    `json:"target_id" bson:"target_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// SessionMessage representa el mensaje enriquecido para ser enviado al bus de datos
type SessionMessage struct {
	RecipientID string          `json:"recipient_id" bson:"recipient_id"` // ID del destinatario
	Browser     string          `json:"browser" bson:"browser"`           // Navegador utilizado por el usuario
	IP          string          `json:"ip" bson:"ip"`                     // Dirección IP del usuario
	SO          string          `json:"so" bson:"so"`                     // Sistema operativo del cliente
	Message     MessageReceived `json:"message" bson:"message"`           // Mensaje enviado por el usuario
	TraceID     string          `json:"trace_id" bson:"trace_id"`         // ID de la traza
}

// --- Satisfacer la interfaz eventbus.Event del módulo base ---

func (m SessionMessage) Name() string { return "session.message" }

func (m SessionMessage) Version() string { return "1.0.0" }

func (m SessionMessage) OccurredAt() time.Time { return m.Message.CreatedAt }

func (m SessionMessage) Payload() any { return m }

func (m SessionMessage) Metadata() eventbus.Metadata {
	return eventbus.Metadata{
		"trace_id":     m.TraceID,
		"recipient_id": m.RecipientID,
	}
}

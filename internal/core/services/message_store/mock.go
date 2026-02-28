package message_store

import (
	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/entities"
	"github.com/jalocompany/orbit/internal/core/ports"
)

type messageStoreMock struct{}

// NewMessageStoreMock instancia el servicio mock de almacenamiento de mensajes.
func NewMessageStoreMock() ports.MessageStoreService {
	return &messageStoreMock{}
}

func (s *messageStoreMock) SaveMessage(ctx *customctx.CustomContext, msg entities.SessionMessage) error {
	ctx.Logger().Info("[MOCK MessageStore] Mensaje de %s para %s guardado exitosamente: %s",
		msg.Message.SenderID, msg.RecipientID, msg.Message.Text)
	return nil
}

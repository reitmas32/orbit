package ports

import (
	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/entities"
)

// ChatService define el puerto para el servicio de chat, enfocado en el bus de mensajes.
type ChatService interface {
	PublishMessage(ctx *customctx.CustomContext, msg entities.SessionMessage) error
}

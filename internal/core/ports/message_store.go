package ports

import (
	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/entities"
)

// MessageStoreService define la interfaz para el almacenamiento persistente de mensajes.
type MessageStoreService interface {
	SaveMessage(ctx *customctx.CustomContext, msg entities.SessionMessage) error
}

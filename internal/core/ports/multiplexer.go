package ports

import (
	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/entities"
)

// MultiplexerService define la interfaz para el servicio de multiplexado de mensajes.
type MultiplexerService interface {
	Multiplex(ctx *customctx.CustomContext, msg entities.SessionMessage) error
}

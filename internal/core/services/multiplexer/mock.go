package multiplexer

import (
	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/entities"
	"github.com/jalocompany/orbit/internal/core/ports"
)

type multiplexerMock struct{}

// NewMultiplexerMock instancia el servicio mock de multiplexado.
func NewMultiplexerMock() ports.MultiplexerService {
	return &multiplexerMock{}
}

func (s *multiplexerMock) Multiplex(ctx *customctx.CustomContext, msg entities.SessionMessage) error {
	ctx.Logger().Info("[MOCK Multiplexer] Mensaje de grupo hacia %s multiplexado exitosamente", msg.Message.TargetID)
	return nil
}

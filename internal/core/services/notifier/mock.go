package notifier

import (
	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/ports"
)

type userNotifierMock struct{}

// NewUserNotifierMock instancia el servicio mock de notificaciones
func NewUserNotifierMock() ports.UserNotifierService {
	return &userNotifierMock{}
}

func (s *userNotifierMock) NotifyConnect(ctx *customctx.CustomContext, userID string) error {
	ctx.Logger().Info("[MOCK Notifier] User %s connected", userID)
	return nil
}

func (s *userNotifierMock) NotifyDisconnect(ctx *customctx.CustomContext, userID string) error {
	ctx.Logger().Info("[MOCK Notifier] User %s disconnected", userID)
	return nil
}

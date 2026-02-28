package ports

import "github.com/fundathyon/base/core/customctx"

// UserNotifierService define la interfaz para el servicio de notificaciones de conexión de usuarios.
type UserNotifierService interface {
	NotifyConnect(ctx *customctx.CustomContext, userID string) error
	NotifyDisconnect(ctx *customctx.CustomContext, userID string) error
}

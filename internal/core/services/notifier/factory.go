package notifier

import (
	"strings"

	"github.com/jalocompany/orbit/internal/core/ports"
	"github.com/jalocompany/orbit/internal/core/settings"
)

// UserNotifierFactory emplea el patrón factory para retornar la implementación mock o la real.
func UserNotifierFactory() ports.UserNotifierService {
	// Si no hay URL configurada o la URL es el tag 'mock', retornamos el mock
	if settings.Settings.CORE_API_URL == "" || strings.ToLower(settings.Settings.CORE_API_URL) == "mock" {
		return NewUserNotifierMock()
	}

	// Si sí está configurada, usar el servicio real HTTP
	return NewUserNotifierHTTP(settings.Settings.CORE_API_URL, settings.Settings.CORE_API_KEY)
}

package multiplexer

import (
	"github.com/jalocompany/orbit/internal/core/ports"
	"github.com/jalocompany/orbit/internal/core/settings"
)

// MultiplexerFactory crea e inicializa el servicio de multiplexado acorde a la configuración.
func MultiplexerFactory() ports.MultiplexerService {
	if settings.Settings.MULTIPLEXER_API_URL == "mock" {
		return NewMultiplexerMock()
	}
	return NewMultiplexerHTTP(settings.Settings.MULTIPLEXER_API_URL, settings.Settings.MULTIPLEXER_API_KEY)
}

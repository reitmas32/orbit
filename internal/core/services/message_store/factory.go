package message_store

import (
	"github.com/jalocompany/orbit/internal/core/ports"
	"github.com/jalocompany/orbit/internal/core/settings"
)

// MessageStoreFactory crea e inicializa el servicio de almacenamiento de mensajes acorde a la configuración.
func MessageStoreFactory() ports.MessageStoreService {
	if settings.Settings.CORE_API_URL == "mock" {
		return NewMessageStoreMock()
	}
	return NewMessageStoreHTTP(settings.Settings.CORE_API_URL, settings.Settings.CORE_API_KEY)
}

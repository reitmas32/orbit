package message_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/entities"
	"github.com/jalocompany/orbit/internal/core/ports"
)

type messageStoreHTTP struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewMessageStoreHTTP instancia el servicio real de almacenamiento de mensajes vía CORE API.
func NewMessageStoreHTTP(baseURL, apiKey string) ports.MessageStoreService {
	return &messageStoreHTTP{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

func (s *messageStoreHTTP) SaveMessage(ctx *customctx.CustomContext, msg entities.SessionMessage) error {
	ctx.Logger().Info("[HTTP MessageStore] Persistiendo mensaje de usuario %s hacia CORE API", msg.Message.SenderID)

	body, _ := json.Marshal(msg)

	// Endpoint propuesto para almacenar mensajes
	url := fmt.Sprintf("%s/v1/messages", s.baseURL)

	req, err := http.NewRequestWithContext(ctx.Context(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		ctx.Logger().Error("[HTTP MessageStore] Failed to create request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
		req.Header.Set("x-api-key", s.apiKey)
	}

	// Traspasar trace ID si existe
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		req.Header.Set("x-trace-id", traceID)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		ctx.Logger().Error("[HTTP MessageStore] Request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		ctx.Logger().Error("[HTTP MessageStore] Request returned failing status code %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	ctx.Logger().Info("[HTTP MessageStore] Mensaje persistido exitosamente en Core")
	return nil
}

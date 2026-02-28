package multiplexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/entities"
	"github.com/jalocompany/orbit/internal/core/ports"
)

type multiplexerHTTP struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewMultiplexerHTTP instancia el servicio real de multiplexado vía CORE/MULTIPLEXER API.
func NewMultiplexerHTTP(baseURL, apiKey string) ports.MultiplexerService {
	return &multiplexerHTTP{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

func (s *multiplexerHTTP) Multiplex(ctx *customctx.CustomContext, msg entities.SessionMessage) error {
	ctx.Logger().Info("[HTTP Multiplexer] Multiplexando mensaje de grupo %s hacia Multiplexer API", msg.Message.TargetID)

	body, _ := json.Marshal(msg)

	// Endpoint propuesto para multiplexar mensajes a un grupo
	url := fmt.Sprintf("%s/v1/multiplex", s.baseURL)

	req, err := http.NewRequestWithContext(ctx.Context(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		ctx.Logger().Error("[HTTP Multiplexer] Failed to create request: %v", err)
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
		ctx.Logger().Error("[HTTP Multiplexer] Request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		ctx.Logger().Error("[HTTP Multiplexer] Request returned failing status code %d", resp.StatusCode)
		return fmt.Errorf("multiplexer request failed: status code %d", resp.StatusCode)
	}

	ctx.Logger().Info("[HTTP Multiplexer] Mensaje multiplexado exitosamente")
	return nil
}

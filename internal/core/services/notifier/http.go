package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/ports"
)

type userNotifierHTTP struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewUserNotifierHTTP instancia el servicio real de notificaciones vía HTTP
func NewUserNotifierHTTP(baseURL, apiKey string) ports.UserNotifierService {
	return &userNotifierHTTP{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

func (s *userNotifierHTTP) NotifyConnect(ctx *customctx.CustomContext, userID string) error {
	return s.doRequest(ctx, "connect", userID)
}

func (s *userNotifierHTTP) NotifyDisconnect(ctx *customctx.CustomContext, userID string) error {
	return s.doRequest(ctx, "disconnect", userID)
}

func (s *userNotifierHTTP) doRequest(ctx *customctx.CustomContext, action, userID string) error {
	ctx.Logger().Info("[HTTP Notifier] Notifying %s for user %s via Core API", action, userID)

	payload := map[string]string{
		"user_id": userID,
		"action":  action,
	}
	body, _ := json.Marshal(payload)

	// Endpoint genérico propuesto, se concatena a la URL del core
	url := fmt.Sprintf("%s/v1/users/%s", s.baseURL, action)

	req, err := http.NewRequestWithContext(ctx.Context(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		ctx.Logger().Error("[HTTP Notifier] Failed to create request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
		// Alternativamente si esperan header personalizado para el apikey:
		req.Header.Set("x-api-key", s.apiKey)
	}

	// Traspasar trace ID si existe
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		req.Header.Set("x-trace-id", traceID)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		ctx.Logger().Error("[HTTP Notifier] Request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		ctx.Logger().Error("[HTTP Notifier] Request returned failing status code %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	ctx.Logger().Info("[HTTP Notifier] Successfully notified %s for user %s", action, userID)
	return nil
}

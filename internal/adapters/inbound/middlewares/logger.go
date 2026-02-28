package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fundathyon/base/core/customctx"
	"github.com/fundathyon/base/observability/logger/loguru"
	"github.com/fundathyon/base/observability/logger/loguru/fields"
	"github.com/jalocompany/orbit/internal/core/settings"
)

// WebSocketLoggerMiddleware contextualiza el logger.
func WebSocketLoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Crear el CustomContext instanciado a partir del actual
		ctx := customctx.New(r.Context())

		// Extraer headers de trazabilidad
		traceID := r.Header.Get("X-Trace-ID")
		callerID := r.Header.Get("X-Caller-ID")

		// Instanciar un nuevo Logger usando WSFields
		wsFields := &fields.WSFields{
			Path:     r.URL.Path,
			TraceID:  traceID,
			ClientID: callerID, // Asignamos X-Caller-ID a ClientID
			ServerID: fmt.Sprintf("%s", settings.Settings.SERVER_ID),
		}

		logr := loguru.NewLogger(wsFields)

		// Guardar al logger dentro del wrapper the customctx y propagarlo a la request
		ctxWithLogger := ctx.WithLogger(logr)

		log.Printf("[Middleware: Logger] Logger Loguro inicializado para TraceID: %s, CallerID: %s\n", traceID, callerID)

		next(w, r.WithContext(ctxWithLogger))
	}
}

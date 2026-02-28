package middlewares

import (
	"net/http"
	"strings"

	"github.com/fundathyon/base/core/customctx"
	"github.com/jalocompany/orbit/internal/core/auth"
	"github.com/jalocompany/orbit/internal/core/settings"
)

// JWTAuthMiddleware verifica la autenticación JWT para conexiones
func JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, ok := r.Context().(*customctx.CustomContext)
		if !ok {
			// Si no hay context, lo creamos
			ctx = customctx.New(r.Context())
		}

		var tokenString string

		// 1. Intentar extraer de la cabecera Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = parts[1]
			}
		}

		// 2. Fallback a query param 'access_token' (común en WebSockets ya que no todos los clientes soportan headers)
		if tokenString == "" {
			tokenString = r.URL.Query().Get("access_token")
		}

		if tokenString == "" {
			ctx.Logger().Error("[Middleware: JWTAuth] Token is missing in Authorization header or query params")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &auth.JWTClaims{}

		// Usar nuestro JWTController para validar firmemente el token
		jwtCtrl := auth.JWTController{PublicKey: settings.Settings.JWT_PUBLIC_KEY}
		err := jwtCtrl.ValidateToken(ctx, tokenString, claims)

		if err != nil {
			ctx.Logger().Error("[Middleware: JWTAuth] Token validation failed: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx.Logger().Info("[Middleware: JWTAuth] Token successfully validated for UserID: %s", claims.UserID)

		// Inyectar claims en el CustomContext para que los handlers/servicios tengan acceso
		ctxWithClaims := ctx.WithValue("user_id", claims.UserID)
		ctxWithClaims = ctxWithClaims.WithValue("claims", claims)

		ctxWithClaims.Logger().Info("[Middleware: JWTAuth] Claims: %s", claims.String())

		next(w, r.WithContext(ctxWithClaims))
	}
}

package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/fundathyon/base/core/customctx"
	"github.com/golang-jwt/jwt/v4"
)

// JWTController gestiona generación y validación de JWTs usando RSA (RS256).
// PrivateKey y PublicKey deben contener las claves en formato PEM.
type JWTController struct {
	PublicKey string
}

func normalizePEM(pemStr string) string {
	const headerMarker = "-----BEGIN PRIVATE KEY-----"
	const footerMarker = "-----END PRIVATE KEY-----"

	// Si ya tiene saltos de línea, asumimos que está bien.
	if strings.Contains(pemStr, "\n") {
		return pemStr
	}

	// Buscamos los índices exactos del encabezado y pie.
	beginIdx := strings.Index(pemStr, headerMarker)
	endIdx := strings.Index(pemStr, footerMarker)
	if beginIdx == -1 || endIdx == -1 {
		// No encontramos los marcadores completos: devolvemos tal cual.
		return pemStr
	}

	// La parte entre headerMarker y footerMarker (excluyendo los marcadores).
	// Nota: +len(headerMarker) avanza justo después de "-----END PRIVATE KEY-----"
	startOfBody := beginIdx + len(headerMarker)
	bodyAndExtras := pemStr[startOfBody:endIdx]

	// Eliminamos todos los espacios en blanco de la porción base64.
	// Es probable que en un solo string venga así: " MIIEvQIBAD…abC...  " (con espacios intermedios).
	onlyBase64 := strings.ReplaceAll(bodyAndExtras, " ", "")
	onlyBase64 = strings.ReplaceAll(onlyBase64, "\t", "")
	onlyBase64 = strings.ReplaceAll(onlyBase64, "\r", "")
	onlyBase64 = strings.ReplaceAll(onlyBase64, "\n", "")

	// Reensamblamos: header + newline + body En líneas de 64 chars + newline + footer + newline
	var formattedBody strings.Builder
	for i := 0; i < len(onlyBase64); i += 64 {
		end := i + 64
		if end > len(onlyBase64) {
			end = len(onlyBase64)
		}
		formattedBody.WriteString(onlyBase64[i:end] + "\n")
	}

	return fmt.Sprintf("%s\n%s%s\n", headerMarker, formattedBody.String(), footerMarker)
}

func normalizePEMPublicKey(pemStr string) string {
	const headerMarker = "-----BEGIN PUBLIC KEY-----"
	const footerMarker = "-----END PUBLIC KEY-----"

	// Si ya tiene saltos de línea, asumimos que está bien.
	if strings.Contains(pemStr, "\n") {
		return pemStr
	}

	// Buscamos los índices exactos del encabezado y pie.
	beginIdx := strings.Index(pemStr, headerMarker)
	endIdx := strings.Index(pemStr, footerMarker)
	if beginIdx == -1 || endIdx == -1 {
		// No encontramos los marcadores completos: devolvemos tal cual.
		return pemStr
	}

	// La parte entre headerMarker y footerMarker (excluyendo los marcadores).
	// Nota: +len(headerMarker) avanza justo después de "-----END PRIVATE KEY-----"
	startOfBody := beginIdx + len(headerMarker)
	bodyAndExtras := pemStr[startOfBody:endIdx]

	// Eliminamos todos los espacios en blanco de la porción base64.
	// Es probable que en un solo string venga así: " MIIEvQIBAD…abC...  " (con espacios intermedios).
	onlyBase64 := strings.ReplaceAll(bodyAndExtras, " ", "")
	onlyBase64 = strings.ReplaceAll(onlyBase64, "\t", "")
	onlyBase64 = strings.ReplaceAll(onlyBase64, "\r", "")
	onlyBase64 = strings.ReplaceAll(onlyBase64, "\n", "")

	// Reensamblamos: header + newline + body En líneas de 64 chars + newline + footer + newline
	var formattedBody strings.Builder
	for i := 0; i < len(onlyBase64); i += 64 {
		end := i + 64
		if end > len(onlyBase64) {
			end = len(onlyBase64)
		}
		formattedBody.WriteString(onlyBase64[i:end] + "\n")
	}

	return fmt.Sprintf("%s\n%s%s\n", headerMarker, formattedBody.String(), footerMarker)
}

// parseRSAPrivateKey decodifica un PEM PKCS#1 o PKCS#8 en *rsa.PrivateKey.
func parseRSAPrivateKey(ctx *customctx.CustomContext, pemStr string) (*rsa.PrivateKey, error) {
	entry := ctx.Logger()
	entry.Info("Parsing RSA private key")

	pemStrNormalized := normalizePEM(pemStr)
	//pemStrNormalized := pemStr

	block, _ := pem.Decode([]byte(pemStrNormalized))
	if block == nil {
		entry.Error("Failed to decode PEM private key")
		return nil, errors.New("no se pudo decodificar PEM de clave privada")
	}
	// PKCS#1
	if priv, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		entry.Info("RSA private key parsed successfully")
		return priv, nil
	}
	// PKCS#8
	keyIfc, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		entry.Error("Failed to parse PKCS#8 private key", "error", err)
		return nil, fmt.Errorf("falló parseo PKCS#8: %w", err)
	}
	priv, ok := keyIfc.(*rsa.PrivateKey)
	if !ok {
		entry.Error("PKCS#8 private key is not RSA")
		return nil, errors.New("clave PKCS#8 no es RSA")
	}
	entry.Info("RSA private key parsed successfully")
	return priv, nil
}

// parseRSAPublicKey decodifica un PEM en *rsa.PublicKey.
func parseRSAPublicKey(ctx *customctx.CustomContext, pemStr string) (*rsa.PublicKey, error) {
	entry := ctx.Logger()
	entry.Info("Parsing RSA public key")

	pemStrNormalized := normalizePEMPublicKey(pemStr)
	//pemStrNormalized := pemStr

	block, _ := pem.Decode([]byte(pemStrNormalized))
	if block == nil {
		entry.Error("Failed to decode PEM public key")
		return nil, errors.New("no se pudo decodificar PEM de clave pública")
	}
	pubIfc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		entry.Error("Failed to parse public key", "error", err)
		return nil, fmt.Errorf("falló parseo clave pública: %w", err)
	}
	pub, ok := pubIfc.(*rsa.PublicKey)
	if !ok {
		entry.Error("Public key is not RSA")
		return nil, errors.New("clave pública no es RSA")
	}
	entry.Info("RSA public key parsed successfully")
	return pub, nil
}

// ValidateToken valida un JWT firmado con RS256 y parsea los claims al objeto provisto.
func (j *JWTController) ValidateToken(ctx *customctx.CustomContext, tokenString string, claims jwt.Claims) error {
	entry := ctx.Logger()
	entry.Info("Validating JWT token")

	pubKey, err := parseRSAPublicKey(ctx, j.PublicKey)
	if err != nil {
		entry.Error("Error parsing public key: %v", err)
		return err
	}

	entry.Info("Public key parsed successfully")
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			entry.Error("Unexpected algorithm: %v", token.Method.Alg())
			return nil, fmt.Errorf("algoritmo inesperado: %s", token.Method.Alg())
		}
		return pubKey, nil
	})

	entry.Info("Token parsed successfully")
	// Si hay un error al parsear, comprobamos si se debe a expiración
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
			entry.Warn("Token expired, but allowing connection for now: %v", ve.Inner)
			return nil
		}
		entry.Error("Error parsing token: %v", err)
		return err
	}

	// Si el token es válido, retornamos exitosamente
	if token.Valid {
		entry.Info("Token validated successfully")
		return nil
	}

	entry.Error("Token invalid")
	return errors.New("token inválido")
}

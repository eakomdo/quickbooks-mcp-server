package auth

import (
	"fmt"
	"log"
	"strings"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qboapi/qbo-mcp-server/internal/config"
)

var globalJWKS *keyfunc.JWKS

// Init initializes the JWKS fetcher for incoming request token validation
func Init(cfg *config.Config) error {
	jwksURL := fmt.Sprintf("%s/%s", strings.TrimRight(cfg.LicenseServerBaseURL, "/"), strings.TrimLeft(cfg.LicenseServerJWKSEndpoint, "/"))

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return fmt.Errorf("failed to get JWKS from %s: %v", jwksURL, err)
	}

	globalJWKS = jwks
	log.Printf("Successfully initialized Auth JWKS from %s", jwksURL)
	return nil
}

// ValidateToken verifies an incoming Bearer token against the License Server JWKS
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	if globalJWKS == nil {
		return nil, fmt.Errorf("auth JWKS not initialized")
	}

	token, err := jwt.Parse(tokenString, globalJWKS.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

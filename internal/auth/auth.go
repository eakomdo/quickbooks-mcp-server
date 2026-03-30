package auth

import (
	"fmt"
	"log"
	"strings"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qboapi/qbo-mcp-server/internal/config"
)

var (
	licenseJWKS *keyfunc.JWKS
	accountJWKS *keyfunc.JWKS
)

// Init initializes the JWKS fetchers for incoming request token validation
func Init(cfg *config.Config) error {
	// 1. License Server JWKS
	if cfg.LicenseServerBaseURL != "" {
		jwksURL := fmt.Sprintf("%s/%s", strings.TrimRight(cfg.LicenseServerBaseURL, "/"), strings.TrimLeft(cfg.LicenseServerJWKSEndpoint, "/"))
		jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
		if err != nil {
			log.Printf("Warning: failed to get JWKS from License Server %s: %v", jwksURL, err)
		} else {
			licenseJWKS = jwks
			log.Printf("Successfully initialized License Server JWKS from %s", jwksURL)
		}
	}

	// 2. Account Service JWKS
	if cfg.AccountServiceURL != "" {
		jwksURL := fmt.Sprintf("%s/%s", strings.TrimRight(cfg.AccountServiceURL, "/"), strings.TrimLeft(cfg.AccountServiceJWKS, "/"))
		jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
		if err != nil {
			log.Printf("Warning: failed to get JWKS from Account Service %s: %v", jwksURL, err)
		} else {
			accountJWKS = jwks
			log.Printf("Successfully initialized Account Service JWKS from %s", jwksURL)
		}
	}

	if licenseJWKS == nil && accountJWKS == nil {
		return fmt.Errorf("failed to initialize any JWKS source")
	}

	return nil
}

// ValidateToken verifies an incoming Bearer token against available JWKS sources
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	var lastErr error

	// Try License Server JWKS if available
	if licenseJWKS != nil {
		token, err := jwt.Parse(tokenString, licenseJWKS.Keyfunc)
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				return claims, nil
			}
		}
		lastErr = err
	}

	// Try Account Service JWKS if available
	if accountJWKS != nil {
		token, err := jwt.Parse(tokenString, accountJWKS.Keyfunc)
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				return claims, nil
			}
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, fmt.Errorf("token validation failed: %v", lastErr)
	}

	return nil, fmt.Errorf("invalid token: no valid JWKS found or token invalid")
}

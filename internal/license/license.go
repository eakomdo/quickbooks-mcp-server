package license

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qboapi/qbo-mcp-server/internal/config"
)

type ActivationPayload struct {
	DeviceID   string         `json:"device_id"`
	LicenseKey string         `json:"license_key"`
	ServiceID  string         `json:"service_id"`
	Metrics    map[string]any `json:"metrics"`
}

type ActivationResponse struct {
	ActivationToken string `json:"activation_token"`
}

func GetDeviceID() string {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && len(i.HardwareAddr) > 0 {
				return fmt.Sprintf("%012x", i.HardwareAddr)
			}
		}
	}
	return "000000000000"
}

func Validate(cfg *config.Config) (bool, jwt.MapClaims) {
	if cfg.LicenseServerBaseURL == "" || cfg.LicenseKey == "" {
		log.Println("License configuration missing.")
		return false, nil
	}

	deviceID := GetDeviceID()
	if cfg.DeviceID != "" {
		deviceID = cfg.DeviceID
	}

	activationURL := fmt.Sprintf("%s/%s", strings.TrimRight(cfg.LicenseServerBaseURL, "/"), strings.TrimLeft(cfg.LicenseServerActivationPath, "/"))
	jwksURL := fmt.Sprintf("%s/%s", strings.TrimRight(cfg.LicenseServerBaseURL, "/"), strings.TrimLeft(cfg.LicenseServerJWKSEndpoint, "/"))

	payload := ActivationPayload{
		DeviceID:   deviceID,
		LicenseKey: cfg.LicenseKey,
		ServiceID:  cfg.ServiceID,
		Metrics: map[string]any{
			"example_metric": 21,
		},
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(activationURL, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("License activation error: %v", err)
		return false, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("License activation failed with status: %d", resp.StatusCode)
		return false, nil
	}

	var actResp ActivationResponse
	if err := json.NewDecoder(resp.Body).Decode(&actResp); err != nil {
		log.Printf("Invalid activation response: %v", err)
		return false, nil
	}

	if actResp.ActivationToken == "" {
		log.Println("No activation token returned")
		return false, nil
	}

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		log.Printf("Failed to get JWKS from %s: %v", jwksURL, err)
		return false, nil
	}

	token, err := jwt.Parse(actResp.ActivationToken, jwks.Keyfunc)
	if err != nil {
		log.Printf("Failed to parse/validate license token: %v", err)
		return false, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sid, ok := claims["service_id"].(string); !ok || sid != cfg.ServiceID {
			log.Printf("License service_id mismatch: got %v, expected %s", claims["service_id"], cfg.ServiceID)
			return false, nil
		}
		if did, ok := claims["device_id"].(string); !ok || did != deviceID {
			log.Printf("License device_id mismatch: got %v, expected %s", claims["device_id"], deviceID)
			return false, nil
		}
		return true, claims
	}

	log.Println("Invalid license claims")
	return false, nil
}

func Watcher(cfg *config.Config, intervalSeconds int, maxFailures int) {
	if intervalSeconds <= 0 {
		intervalSeconds = 86400
	}
	if maxFailures <= 0 {
		maxFailures = 14
	}
	failureCount := 0

	for {
		time.Sleep(time.Duration(intervalSeconds) * time.Second)

		valid, _ := Validate(cfg)
		if valid {
			failureCount = 0
			log.Println("License check passed!!!!!!")
		} else {
			failureCount++
			log.Printf("License check failed %d times.", failureCount)
			if failureCount >= maxFailures {
				log.Println("Max failures reached. Shutting down...")
				os.Exit(1)
			}
		}
	}
}

package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds environment-driven settings (mirrors the TypeScript server).
type Config struct {
	QuickBooksClientID     string
	QuickBooksClientSecret string
	QuickBooksRefreshToken string
	QuickBooksRealmID      string
	QuickBooksEnvironment  string // "sandbox" or "production"

	PlatformIntURL string

	EnforceAuth bool

	ServiceID            string
	SecretKey            string // v1 JWT HMAC per client_id
	AccountServiceURL    string
	AccountServiceJWKS   string
	AccountServiceJWKSTTL int // seconds

	UsageReportEndpoint string

	LicenseServerBaseURL        string
	LicenseServerJWKSEndpoint   string
	LicenseServerActivationPath string
	LicenseKey                  string

	HTTPPort   string
	HTTPPath   string
	DeviceID   string
	RedirectURI string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		QuickBooksClientID:     os.Getenv("QUICKBOOKS_CLIENT_ID"),
		QuickBooksClientSecret: os.Getenv("QUICKBOOKS_CLIENT_SECRET"),
		QuickBooksRefreshToken: os.Getenv("QUICKBOOKS_REFRESH_TOKEN"),
		QuickBooksRealmID:      os.Getenv("QUICKBOOKS_REALM_ID"),
		QuickBooksEnvironment:  strings.ToLower(strings.TrimSpace(getenvDefault("QUICKBOOKS_ENVIRONMENT", "sandbox"))),

		PlatformIntURL: strings.TrimRight(os.Getenv("PLATFORM_INT_URL"), "/"),

		EnforceAuth: strings.EqualFold(os.Getenv("ENFORCE_AUTH"), "true"),

		ServiceID:            getenvDefault("SERVICE_ID", "quickbooks-mcp"),
		SecretKey:            os.Getenv("SECRET_KEY"),
		AccountServiceURL:    strings.TrimRight(os.Getenv("ACCOUNT_SERVICE_URL"), "/"),
		AccountServiceJWKS:  getenvDefault("ACCOUNT_SERVICE_JWKS_ENDPOINT", "/.well-known/jwks.json"),
		AccountServiceJWKSTTL: atoiDefault("ACCOUNT_SERVICE_JWKS_CACHE_TTL", 600),

		UsageReportEndpoint: os.Getenv("USAGE_REPORT_ENDPOINT"),

		LicenseServerBaseURL:        strings.TrimRight(os.Getenv("LICENSE_SERVER_BASE_URL"), "/"),
		LicenseServerJWKSEndpoint:   getenvDefault("LICENSE_SERVER_JWKS_ENDPOINT", "/.well-known/jwks.json"),
		LicenseServerActivationPath: getenvDefault("LICENSE_SERVER_ACTIVATION_ENDPOINT", "/activate"),
		LicenseKey:                  strings.TrimSpace(os.Getenv("LICENSE_KEY")),

		HTTPPort:   os.Getenv("PORT"),
		HTTPPath:   getenvDefault("MCP_HTTP_PATH", "/mcp"),
		DeviceID:   os.Getenv("DEVICE_ID"),
		RedirectURI: getenvDefault("QUICKBOOKS_REDIRECT_URI", "http://localhost:8000/callback"),
	}
	return cfg, nil
}

func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func atoiDefault(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func (c *Config) UsePlatformIntegration() bool {
	return c.PlatformIntURL != ""
}

func (c *Config) QuickBooksAPIBase() string {
	if c.QuickBooksEnvironment == "production" {
		return "https://quickbooks.api.intuit.com/v3/company/"
	}
	return "https://sandbox-quickbooks.api.intuit.com/v3/company/"
}

func (c *Config) OAuthTokenURL() string {
	return "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer"
}

package qbo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/qboapi/qbo-mcp-server/internal/config"
)

const minorVersion = "75"

type ctxKey int

const userBearerCtxKey ctxKey = 1

// WithUserBearer stores the end-user Bearer token (platform mode) on the context.
func WithUserBearer(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, userBearerCtxKey, strings.TrimSpace(token))
}

func userBearerFrom(ctx context.Context) string {
	v, _ := ctx.Value(userBearerCtxKey).(string)
	return v
}

type QBProject struct {
	ID           string `json:"id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	Environment  string `json:"environment"`
	RealmID      string `json:"realmid"`
	RedirectURI  string `json:"redirect_uri"`
}

// Client performs QuickBooks Online v3 REST calls (node-quickbooks compatible paths).
type Client struct {
	cfg *config.Config

	mu sync.Mutex

	// Effective credentials (from env or platform)
	clientID     string
	clientSecret string
	refreshToken string
	realmID      string
	projectID    string

	accessToken string
	expiresAt   time.Time

	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		clientID:     cfg.QuickBooksClientID,
		clientSecret: cfg.QuickBooksClientSecret,
		refreshToken: cfg.QuickBooksRefreshToken,
		realmID:      cfg.QuickBooksRealmID,
	}
}

func (c *Client) loadPlatformProject(ctx context.Context) error {
	if c.cfg.PlatformIntURL == "" {
		return fmt.Errorf("PLATFORM_INT_URL not configured")
	}
	tok := userBearerFrom(ctx)
	if tok == "" {
		return fmt.Errorf("no Authorization token found; provide a Bearer token in the request")
	}
	u := c.cfg.PlatformIntURL + "/api/v1/quickbooks_projects?limit=1"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("platform integration unreachable: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to fetch QB project from platform: %s %s", resp.Status, string(body))
	}
	var data struct {
		Items []QBProject `json:"items"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("decode platform response: %w", err)
	}
	if len(data.Items) == 0 {
		return fmt.Errorf("no QuickBooks project found in platform integration service")
	}
	p := data.Items[0]
	c.projectID = p.ID
	c.clientID = p.ClientID
	c.clientSecret = p.ClientSecret
	c.refreshToken = p.RefreshToken
	if p.RealmID != "" {
		c.realmID = p.RealmID
	}
	return nil
}

func (c *Client) syncRefreshTokenToPlatform(ctx context.Context, newRefresh, userBearer string) {
	if c.cfg.PlatformIntURL == "" || c.projectID == "" || userBearer == "" {
		return
	}
	u := fmt.Sprintf("%s/api/v1/quickbooks_projects/%s", c.cfg.PlatformIntURL, c.projectID)
	payload := map[string]string{"refresh_token": newRefresh}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u, bytes.NewReader(b))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+userBearer)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("warning: failed to sync refresh token to platform: %s %s\n", resp.Status, string(body))
	}
}

// Authenticate ensures access token and realm are available (mirrors TS quickbooksClient.authenticate() without browser OAuth).
func (c *Client) Authenticate(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cfg.UsePlatformIntegration() {
		if err := c.loadPlatformProject(ctx); err != nil {
			return err
		}
	} else {
		if c.clientID == "" || c.clientSecret == "" {
			return fmt.Errorf("QUICKBOOKS_CLIENT_ID and QUICKBOOKS_CLIENT_SECRET must be set (or configure PLATFORM_INT_URL)")
		}
	}

	if c.refreshToken == "" || c.realmID == "" {
		return fmt.Errorf("QuickBooks not connected: set QUICKBOOKS_REFRESH_TOKEN and QUICKBOOKS_REALM_ID (or complete OAuth in the Node helper / Intuit playground)")
	}

	if c.accessToken != "" && time.Now().Before(c.expiresAt.Add(-2*time.Minute)) {
		return nil
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", c.refreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.OAuthTokenURL(), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.clientID, c.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("oauth token refresh failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("oauth token refresh failed: %s %s", resp.Status, string(body))
	}

	var tok struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tok); err != nil {
		return fmt.Errorf("decode oauth response: %w", err)
	}
	if tok.AccessToken == "" {
		return fmt.Errorf("oauth response missing access_token")
	}
	c.accessToken = tok.AccessToken
	exp := tok.ExpiresIn
	if exp <= 0 {
		exp = 3600
	}
	c.expiresAt = time.Now().Add(time.Duration(exp) * time.Second)
	if tok.RefreshToken != "" && tok.RefreshToken != c.refreshToken {
		c.refreshToken = tok.RefreshToken
		ut := userBearerFrom(ctx)
		rt := c.refreshToken
		go func() {
			c.syncRefreshTokenToPlatform(context.Background(), rt, ut)
		}()
	}
	return nil
}

func (c *Client) baseURL() string {
	return c.cfg.QuickBooksAPIBase() + c.realmID
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any) (json.RawMessage, error) {
	if err := c.Authenticate(ctx); err != nil {
		return nil, err
	}
	u := c.baseURL() + path
	if query == nil {
		query = url.Values{}
	}
	query.Set("minorversion", minorVersion)
	query.Set("format", "json")
	u = u + "?" + query.Encode()

	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "qbo-mcp-server-go/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var probe map[string]json.RawMessage
	_ = json.Unmarshal(raw, &probe)
	if probe != nil {
		if f, ok := probe["Fault"]; ok && len(f) > 0 {
			return nil, fmt.Errorf("quickbooks fault: %s", string(raw))
		}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("quickbooks http %s: %s", resp.Status, string(raw))
	}
	return raw, nil
}

// Create POST /{entity}
func (c *Client) Create(ctx context.Context, entity string, payload any) (json.RawMessage, error) {
	return c.do(ctx, http.MethodPost, "/"+strings.ToLower(entity), nil, payload)
}

// Read GET /{entity}/{id}
func (c *Client) Read(ctx context.Context, entity, id string) (json.RawMessage, error) {
	path := "/" + strings.ToLower(entity) + "/" + url.PathEscape(id)
	return c.do(ctx, http.MethodGet, path, nil, nil)
}

// Update POST /{entity}?operation=update
func (c *Client) Update(ctx context.Context, entity string, payload any) (json.RawMessage, error) {
	q := url.Values{"operation": {"update"}}
	return c.do(ctx, http.MethodPost, "/"+strings.ToLower(entity), q, payload)
}

// Delete POST /{entity}?operation=delete
func (c *Client) Delete(ctx context.Context, entity string, idOrEntity any) (json.RawMessage, error) {
	var payload any
	switch v := idOrEntity.(type) {
	case string:
		raw, err := c.Read(ctx, entity, v)
		if err != nil {
			return nil, err
		}
		var wrap map[string]json.RawMessage
		if err := json.Unmarshal(raw, &wrap); err != nil {
			return nil, err
		}
		key := CapitalizeEntity(entity)
		inner := wrap[key]
		if len(inner) == 0 {
			return nil, fmt.Errorf("unexpected read response for delete: %s", string(raw))
		}
		if err := json.Unmarshal(inner, &payload); err != nil {
			return nil, err
		}
	default:
		payload = idOrEntity
	}
	q := url.Values{"operation": {"delete"}}
	return c.do(ctx, http.MethodPost, "/"+strings.ToLower(entity), q, payload)
}

// CapitalizeEntity converts e.g. billPayment -> BillPayment (JSON / unwrap key).
func CapitalizeEntity(entity string) string {
	if entity == "" {
		return entity
	}
	r := []rune(entity)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

package maxsattapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

// OAuth2Manager manages OAuth2 authentication with automatic token refresh
// Uses client_credentials grant flow for service account authentication
type OAuth2Manager struct {
	authURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	logger       logger.Logger

	mu          sync.RWMutex
	accessToken string
	expiresAt   time.Time
}

// OAuth2Response represents the OAuth2 token response
type OAuth2Response struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // Seconds
}

// NewOAuth2Manager creates a new OAuth2Manager instance
func NewOAuth2Manager(authURL, clientID, clientSecret string, httpClient *http.Client, log logger.Logger) *OAuth2Manager {
	return &OAuth2Manager{
		authURL:      authURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   httpClient,
		logger:       log,
	}
}

// GetAccessToken returns a valid access token, re-authenticating if necessary
func (o *OAuth2Manager) GetAccessToken(ctx context.Context) (string, error) {
	logger.Debug(ctx, "Started", map[string]any{
		"authURL": o.authURL,
	})

	// Check if current token is still valid with 1 minute buffer
	o.mu.RLock()
	if o.accessToken != "" && time.Now().Add(1*time.Minute).Before(o.expiresAt) {
		token := o.accessToken
		o.mu.RUnlock()
		logger.Debug(ctx, "Using cached access token", map[string]any{
			"expiresAt": o.expiresAt,
		})
		return token, nil
	}
	o.mu.RUnlock()

	// client_credentials grant does not support refresh tokens, re-authenticate on expiry
	logger.Debug(ctx, "Authenticating with client_credentials grant")
	err := o.authenticateWithClientCredentials(ctx)
	if err != nil {
		return "", err
	}

	o.mu.RLock()
	token := o.accessToken
	o.mu.RUnlock()

	logger.Debug(ctx, "Authentication successful")
	return token, nil
}

// authenticateWithClientCredentials performs authentication using client_credentials grant flow
func (o *OAuth2Manager) authenticateWithClientCredentials(ctx context.Context) error {
	logger.Debug(ctx, "Starting client_credentials grant authentication", map[string]any{
		"clientID": o.clientID,
	})

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", o.clientID)
	data.Set("client_secret", o.clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.authURL, strings.NewReader(data.Encode()))
	if err != nil {
		logger.Error(ctx, err, "Failed to create authentication request")
		return errs.InternalServerError(err, "FCMP-01-0001", "Failed to create authentication request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		logger.Error(ctx, err, "Authentication request failed")
		return errs.InternalServerError(err, "FCMP-01-0002", "Authentication request failed")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(ctx, err, "Failed to read authentication response")
		return errs.InternalServerError(err, "FCMP-01-0003", "Failed to read authentication response")
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error(ctx, nil, "Authentication failed with non-200 status", map[string]any{
			"statusCode": resp.StatusCode,
			"response":   string(body),
		})
		return errs.UnauthorizedError("Invalid credentials", "FCMP-01-0401")
	}

	var tokenResp OAuth2Response
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		logger.Error(ctx, err, "Failed to parse authentication response")
		return errs.InternalServerError(err, "FCMP-01-0004", "Failed to parse authentication response")
	}

	o.mu.Lock()
	o.accessToken = tokenResp.AccessToken
	o.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	o.mu.Unlock()

	logger.Info(ctx, "Client credentials authentication successful", map[string]any{
		"expiresIn": tokenResp.ExpiresIn,
		"expiresAt": o.expiresAt,
	})

	return nil
}

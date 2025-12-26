package maxsattapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	http2 "github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/http"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

// MaxsattAPIClient is the base HTTP client for making authenticated requests to the MaxSatt API
type MaxsattAPIClient struct {
	baseURL     string
	httpClient  http2.Client
	authManager *OAuth2Manager
	logger      logger.Logger
	retryConfig RetryConfig
}

// RetryConfig configures retry behavior for failed requests
type RetryConfig struct {
	MaxAttempts  int           // Maximum number of retry attempts
	InitialDelay time.Duration // Initial delay before first retry
	MaxDelay     time.Duration // Maximum delay between retries
	Multiplier   float64       // Backoff multiplier
}

// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     8 * time.Second,
		Multiplier:   2.0,
	}
}

// APIErrorDetail represents a single error detail entry
type APIErrorDetail struct {
	Attribute string   `json:"attribute"`
	Messages  []string `json:"messages"`
}

// APIError represents an error response from the MaxSatt API
type APIError struct {
	Code         string           `json:"code"`
	Description  string           `json:"description"`
	ErrorDetails []APIErrorDetail `json:"error_details,omitempty"`
	ID           string           `json:"id,omitempty"`
}

// APIErrorResponse represents the error response structure
type APIErrorResponse struct {
	Error APIError `json:"error"`
}

// NewMaxsattAPIClient creates a new MaxSatt API client
func NewMaxsattAPIClient(
	baseURL string,
	authURL string,
	clientID string,
	clientSecret string,
	httpClient http2.Client,
	log logger.Logger,
) *MaxsattAPIClient {
	authHTTPClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &MaxsattAPIClient{
		baseURL:     baseURL,
		httpClient:  httpClient,
		authManager: NewOAuth2Manager(authURL, clientID, clientSecret, authHTTPClient, log),
		logger:      log,
		retryConfig: DefaultRetryConfig(),
	}
}

// DoRequest executes an HTTP request with authentication, retry logic, and error handling
func (c *MaxsattAPIClient) DoRequest(
	ctx context.Context,
	method string,
	path string,
	requestBody interface{},
	responseBody interface{},
) error {
	logger.Debug(ctx, "Starting API request", map[string]any{
		"method": method,
		"path":   path,
	})

	url := c.baseURL + path

	var attempt int
	var lastErr error

	for attempt = 0; attempt < c.retryConfig.MaxAttempts; attempt++ {
		if attempt > 0 {
			delay := c.calculateBackoff(attempt)
			logger.Info(ctx, "Retrying request after delay", map[string]any{
				"attempt": attempt + 1,
				"delay":   delay.String(),
			})
			time.Sleep(delay)
		}

		err := c.executeRequest(ctx, method, url, requestBody, responseBody)
		if err == nil {
			logger.Debug(ctx, "API request successful", map[string]any{
				"method":  method,
				"path":    path,
				"attempt": attempt + 1,
			})
			return nil
		}

		lastErr = err

		if !c.isRetryable(err) {
			logger.Info(ctx, "Error is not retryable, stopping", map[string]any{
				"error": err.Error(),
			})
			break
		}

		logger.Warn(ctx, err, "Request failed, will retry", map[string]any{
			"attempt":     attempt + 1,
			"maxAttempts": c.retryConfig.MaxAttempts,
		})
	}

	logger.Error(ctx, lastErr, "API request failed after all retry attempts", map[string]any{
		"method":   method,
		"path":     path,
		"attempts": attempt + 1,
	})

	return lastErr
}

// executeRequest performs a single HTTP request attempt
func (c *MaxsattAPIClient) executeRequest(
	ctx context.Context,
	method string,
	url string,
	requestBody interface{},
	responseBody interface{},
) error {
	accessToken, err := c.authManager.GetAccessToken(ctx)
	if err != nil {
		logger.Error(ctx, err, "Failed to get access token")
		return err
	}

	var bodyReader io.Reader
	if requestBody != nil {
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			logger.Error(ctx, err, "Failed to marshal request body")
			return errs.InternalServerError(err, "FCMP-02-0001", "Failed to marshal request body")
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		logger.Error(ctx, err, "Failed to create HTTP request")
		return errs.InternalServerError(err, "FCMP-02-0002", "Failed to create HTTP request")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(ctx, req)
	if err != nil {
		logger.Error(ctx, err, "HTTP request failed")
		return errs.InternalServerError(err, "FCMP-02-0003", "HTTP request failed")
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(ctx, err, "Failed to read response body")
		return errs.InternalServerError(err, "FCMP-02-0004", "Failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleErrorResponse(ctx, resp.StatusCode, respBodyBytes)
	}

	if responseBody != nil && len(respBodyBytes) > 0 {
		if err := json.Unmarshal(respBodyBytes, responseBody); err != nil {
			logger.Error(ctx, err, "Failed to unmarshal response body", map[string]any{
				"response": string(respBodyBytes),
			})
			return errs.InternalServerError(err, "FCMP-02-0005", "Failed to parse response body")
		}
	}

	return nil
}

// handleErrorResponse processes API error responses and returns appropriate errors
func (c *MaxsattAPIClient) handleErrorResponse(ctx context.Context, statusCode int, body []byte) error {
	var apiErr APIErrorResponse
	if err := json.Unmarshal(body, &apiErr); err != nil {
		logger.Error(ctx, nil, "API returned error status", map[string]any{
			"statusCode": statusCode,
			"response":   string(body),
		})
		return fmt.Errorf("API error (status %d): %s", statusCode, string(body))
	}

	var errorDetailsStr string
	if len(apiErr.Error.ErrorDetails) > 0 {
		var details []string
		for _, detail := range apiErr.Error.ErrorDetails {
			for _, msg := range detail.Messages {
				details = append(details, fmt.Sprintf("[%s] %s", detail.Attribute, msg))
			}
		}
		errorDetailsStr = fmt.Sprintf(": %s", details[0])
		if len(details) > 1 {
			for _, d := range details[1:] {
				errorDetailsStr += fmt.Sprintf(", %s", d)
			}
		}
	}

	logger.Error(ctx, nil, "API returned structured error", map[string]any{
		"statusCode":   statusCode,
		"errorCode":    apiErr.Error.Code,
		"description":  apiErr.Error.Description,
		"errorDetails": apiErr.Error.ErrorDetails,
		"errorID":      apiErr.Error.ID,
	})

	errorMessage := apiErr.Error.Description
	if errorDetailsStr != "" {
		errorMessage = apiErr.Error.Description + errorDetailsStr
	}

	switch statusCode {
	case http.StatusBadRequest:
		return errs.BadRequestError(errorMessage, apiErr.Error.Code)
	case http.StatusUnauthorized:
		return errs.UnauthorizedError(fmt.Sprintf("unauthorized: %s", errorMessage), apiErr.Error.Code)
	case http.StatusForbidden:
		clientErr := errs.ClientError(errorMessage, apiErr.Error.Code, http.StatusForbidden)
		return fmt.Errorf("forbidden: %w", clientErr)
	case http.StatusNotFound:
		return errs.NotFoundError(errorMessage, apiErr.Error.Code)
	case http.StatusConflict:
		return errs.ClientError(errorMessage, apiErr.Error.Code, http.StatusConflict)
	case http.StatusServiceUnavailable:
		return errs.InternalServerError(
			fmt.Errorf("service unavailable: %s", errorMessage),
			"FCMP-02-0503",
			errorMessage,
		)
	default:
		return errs.InternalServerError(
			fmt.Errorf("%s: %s", apiErr.Error.Code, errorMessage),
			"FCMP-02-0500",
			errorMessage,
		)
	}
}

// isRetryable determines if an error is retryable.
// Server errors (5xx) are retried. Client errors (4xx) including 401 are not retried
// since OAuth2Manager already handles token refresh automatically before each request.
// If we still get a 401 here, it means credentials are invalid or token was revoked,
// and retrying would not help.
func (c *MaxsattAPIClient) isRetryable(err error) bool {
	if baseErr, ok := err.(errs.BaseError); ok {
		statusCode := baseErr.StatusCode()

		if statusCode >= 500 && statusCode < 600 {
			return true
		}

		return false
	}

	// Retry network errors (no HTTP response received)
	return true
}

// calculateBackoff calculates the backoff delay for retry attempts
func (c *MaxsattAPIClient) calculateBackoff(attempt int) time.Duration {
	delay := float64(c.retryConfig.InitialDelay) * math.Pow(c.retryConfig.Multiplier, float64(attempt-1))

	if delay > float64(c.retryConfig.MaxDelay) {
		delay = float64(c.retryConfig.MaxDelay)
	}

	return time.Duration(delay)
}

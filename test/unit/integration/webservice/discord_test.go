package webservice_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/webservice"
)

type mockHTTPClient struct {
	doFunc func(ctx context.Context, req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return m.doFunc(ctx, req)
}

func TestNewDiscordNotifier_EmptyWebhookURL_ReturnsNil(t *testing.T) {
	notifier := webservice.NewDiscordNotifier(nil, "")
	if notifier != nil {
		t.Error("NewDiscordNotifier should return nil when webhookURL is empty")
	}
}

func TestNewDiscordNotifier_WithWebhookURL_ReturnsNotifier(t *testing.T) {
	mockClient := &mockHTTPClient{}
	notifier := webservice.NewDiscordNotifier(mockClient, "https://discord.com/api/webhooks/test")
	if notifier == nil {
		t.Error("NewDiscordNotifier should return notifier when webhookURL is provided")
	}
}

func TestDiscordNotifier_Notify_Success(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	notifier := webservice.NewDiscordNotifier(mockClient, "https://discord.com/api/webhooks/test")
	if notifier == nil {
		t.Fatal("NewDiscordNotifier should not return nil")
	}

	payload := adapter.NotificationPayload{
		Level:         adapter.NotificationLevelError,
		Title:         "Test Error",
		Message:       "Test message",
		TransactionID: "tx-123",
		ErrorCode:     "ERR-001",
		Details: map[string]string{
			"Status Code": "500",
		},
	}

	err := notifier.Notify(context.Background(), payload)
	if err != nil {
		t.Errorf("Notify should not return error on success: %v", err)
	}
}

func TestDiscordNotifier_Notify_HTTPError(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return nil, io.EOF
		},
	}

	notifier := webservice.NewDiscordNotifier(mockClient, "https://discord.com/api/webhooks/test")
	if notifier == nil {
		t.Fatal("NewDiscordNotifier should not return nil")
	}

	payload := adapter.NotificationPayload{
		Level:   adapter.NotificationLevelError,
		Title:   "Test Error",
		Message: "Test message",
	}

	err := notifier.Notify(context.Background(), payload)
	if err == nil {
		t.Error("Notify should return error when HTTP request fails")
	}
}

func TestDiscordNotifier_Notify_NonSuccessStatusCode(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	notifier := webservice.NewDiscordNotifier(mockClient, "https://discord.com/api/webhooks/test")
	if notifier == nil {
		t.Fatal("NewDiscordNotifier should not return nil")
	}

	payload := adapter.NotificationPayload{
		Level:   adapter.NotificationLevelError,
		Title:   "Test Error",
		Message: "Test message",
	}

	err := notifier.Notify(context.Background(), payload)
	if err == nil {
		t.Error("Notify should return error when status code is not success")
	}
}

func TestDiscordNotifier_Notify_AllNotificationLevels(t *testing.T) {
	tests := []struct {
		name  string
		level adapter.NotificationLevel
	}{
		{"error level", adapter.NotificationLevelError},
		{"warning level", adapter.NotificationLevelWarning},
		{"info level", adapter.NotificationLevelInfo},
		{"unknown level", adapter.NotificationLevel("unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil
				},
			}

			notifier := webservice.NewDiscordNotifier(mockClient, "https://discord.com/api/webhooks/test")
			if notifier == nil {
				t.Fatal("NewDiscordNotifier should not return nil")
			}

			payload := adapter.NotificationPayload{
				Level:   tt.level,
				Title:   "Test",
				Message: "Test message",
			}

			err := notifier.Notify(context.Background(), payload)
			if err != nil {
				t.Errorf("Notify should not return error: %v", err)
			}
		})
	}
}

func TestDiscordNotifier_Notify_WithAllDetails(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	notifier := webservice.NewDiscordNotifier(mockClient, "https://discord.com/api/webhooks/test")
	if notifier == nil {
		t.Fatal("NewDiscordNotifier should not return nil")
	}

	payload := adapter.NotificationPayload{
		Level:         adapter.NotificationLevelError,
		Title:         "Test Error",
		Message:       "Test message",
		TransactionID: "tx-123",
		ErrorCode:     "ERR-001",
		Details: map[string]string{
			"Status Code":   "500",
			"Retry Status":  "Retriable",
			"Receive Count": "1",
			"Custom Field":  "custom value",
		},
	}

	err := notifier.Notify(context.Background(), payload)
	if err != nil {
		t.Errorf("Notify should not return error: %v", err)
	}
}

func TestDiscordNotifier_Notify_EmptyPayload(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	notifier := webservice.NewDiscordNotifier(mockClient, "https://discord.com/api/webhooks/test")
	if notifier == nil {
		t.Fatal("NewDiscordNotifier should not return nil")
	}

	payload := adapter.NotificationPayload{}

	err := notifier.Notify(context.Background(), payload)
	if err != nil {
		t.Errorf("Notify should not return error with empty payload: %v", err)
	}
}

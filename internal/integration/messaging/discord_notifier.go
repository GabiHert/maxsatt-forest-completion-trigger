package messaging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	pkghttp "github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/http"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type discordNotifier struct {
	httpClient pkghttp.Client
	webhookURL string
}

func NewDiscordNotifier(httpClient pkghttp.Client, webhookURL string) adapter.Notifier {
	return &discordNotifier{
		httpClient: httpClient,
		webhookURL: webhookURL,
	}
}

func (d *discordNotifier) Notify(ctx context.Context, payload adapter.NotificationPayload) error {
	if d.webhookURL == "" {
		logger.Warn(ctx, nil, "Discord webhook URL not configured, skipping notification")
		return nil
	}

	color := d.getColorForLevel(payload.Level)

	fields := make([]map[string]any, 0, len(payload.Details))
	for key, value := range payload.Details {
		fields = append(fields, map[string]any{
			"name":   key,
			"value":  value,
			"inline": true,
		})
	}

	embed := map[string]any{
		"title":       payload.Title,
		"description": payload.Message,
		"color":       color,
		"fields":      fields,
	}

	if payload.ErrorCode != "" {
		embed["footer"] = map[string]string{
			"text": fmt.Sprintf("Error Code: %s | Transaction: %s", payload.ErrorCode, payload.TransactionID),
		}
	}

	discordPayload := map[string]any{
		"embeds": []map[string]any{embed},
	}

	body, err := json.Marshal(discordPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal discord payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create discord request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send discord notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord returned non-success status: %d", resp.StatusCode)
	}

	return nil
}

func (d *discordNotifier) getColorForLevel(level adapter.NotificationLevel) int {
	switch level {
	case adapter.NotificationLevelError:
		return 15158332 // Red
	case adapter.NotificationLevelWarning:
		return 16776960 // Yellow
	case adapter.NotificationLevelInfo:
		return 3447003 // Blue
	default:
		return 9807270 // Gray
	}
}

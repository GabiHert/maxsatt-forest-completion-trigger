package webservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	http2 "github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/http"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

const (
	discordColorRed    = 15548997
	discordColorYellow = 16776960
	discordColorBlue   = 3447003
	discordTimeout     = 5 * time.Second
)

type DiscordMessage struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Color       int          `json:"color"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Footer      *EmbedFooter `json:"footer,omitempty"`
	Timestamp   string       `json:"timestamp,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type EmbedFooter struct {
	Text string `json:"text"`
}

type discordNotifier struct {
	httpClient http2.Client
	webhookURL string
}

func NewDiscordNotifier(httpClient http2.Client, webhookURL string) adapter.Notifier {
	if webhookURL == "" {
		return nil
	}

	return &discordNotifier{
		httpClient: httpClient,
		webhookURL: webhookURL,
	}
}

func (d *discordNotifier) Notify(ctx context.Context, payload adapter.NotificationPayload) error {
	embed := d.buildEmbed(payload)
	message := DiscordMessage{
		Embeds: []DiscordEmbed{embed},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		logger.Warn(ctx, err, "Failed to marshal Discord message")
		return fmt.Errorf("failed to marshal Discord message: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", d.webhookURL, bytes.NewReader(messageBytes))
	if err != nil {
		logger.Warn(ctx, err, "Failed to create Discord webhook request")
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	timeoutCtx, cancel := context.WithTimeout(ctx, discordTimeout)
	defer cancel()

	resp, err := d.httpClient.Do(timeoutCtx, req)
	if err != nil {
		logger.Warn(ctx, err, "Failed to send Discord notification", map[string]any{
			"transactionId": payload.TransactionID,
			"level":         payload.Level,
		})
		return fmt.Errorf("failed to send notification: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Warn(ctx, nil, "Discord webhook returned non-success status", map[string]any{
			"statusCode":    resp.StatusCode,
			"transactionId": payload.TransactionID,
		})
		return fmt.Errorf("Discord webhook returned status %d", resp.StatusCode)
	}

	logger.Debug(ctx, "Discord notification sent successfully", map[string]any{
		"transactionId": payload.TransactionID,
		"level":         payload.Level,
	})

	return nil
}

func (d *discordNotifier) buildEmbed(payload adapter.NotificationPayload) DiscordEmbed {
	embed := DiscordEmbed{
		Title:       payload.Title,
		Description: payload.Message,
		Color:       d.getColorForLevel(payload.Level),
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Footer: &EmbedFooter{
			Text: properties.Properties().Application.ServiceName,
		},
	}

	var fields []EmbedField

	if payload.TransactionID != "" {
		fields = append(fields, EmbedField{
			Name:   "Transaction ID",
			Value:  payload.TransactionID,
			Inline: true,
		})
	}

	if payload.ErrorCode != "" {
		fields = append(fields, EmbedField{
			Name:   "Error Code",
			Value:  payload.ErrorCode,
			Inline: true,
		})
	}

	detailsOrder := []string{"Status Code", "Retry Status", "Receive Count"}
	for _, key := range detailsOrder {
		if value, exists := payload.Details[key]; exists {
			fields = append(fields, EmbedField{
				Name:   key,
				Value:  value,
				Inline: false,
			})
		}
	}

	for key, value := range payload.Details {
		found := false
		for _, orderedKey := range detailsOrder {
			if key == orderedKey {
				found = true
				break
			}
		}
		if !found {
			fields = append(fields, EmbedField{
				Name:   key,
				Value:  value,
				Inline: false,
			})
		}
	}

	if len(fields) > 0 {
		embed.Fields = fields
	}

	return embed
}

func (d *discordNotifier) getColorForLevel(level adapter.NotificationLevel) int {
	switch level {
	case adapter.NotificationLevelError:
		return discordColorRed
	case adapter.NotificationLevelWarning:
		return discordColorYellow
	case adapter.NotificationLevelInfo:
		return discordColorBlue
	default:
		return discordColorRed
	}
}

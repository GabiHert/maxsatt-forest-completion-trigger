package messaging

import (
	"context"
	"fmt"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type failureHandler struct {
	sqsHelper aws.SqsHelperAdapter
	notifier  adapter.Notifier
	dlqURL    string
}

func NewFailureHandler(sqsHelper aws.SqsHelperAdapter, notifier adapter.Notifier, dlqURL string) adapter.FailureHandler {
	return &failureHandler{
		sqsHelper: sqsHelper,
		notifier:  notifier,
		dlqURL:    dlqURL,
	}
}

func (f *failureHandler) Handle(ctx context.Context, err error) error {
	logger.Info(ctx, "Handling failure", map[string]any{
		"error": err.Error(),
	})

	f.sendDiscordNotificationAsync(ctx, err)

	if f.dlqURL == "" {
		logger.Warn(ctx, nil, "DLQ URL not configured, skipping DLQ publish")
		return nil
	}

	message := map[string]any{
		"error":        err.Error(),
		"service_name": "forest-completion-trigger",
	}

	if sqsErr := f.sqsHelper.Send(ctx, message, f.dlqURL); sqsErr != nil {
		logger.Error(ctx, sqsErr, "Failed to send message to DLQ", map[string]any{
			"dlq_url": f.dlqURL,
		})
		return fmt.Errorf("failed to send message to DLQ: %w", sqsErr)
	}

	logger.Info(ctx, "Failure handled and sent to DLQ")

	return nil
}

func (f *failureHandler) sendDiscordNotificationAsync(ctx context.Context, err error) {
	if f.notifier == nil {
		return
	}

	go func() {
		loggerCtx := logger.GetContext(ctx)
		correlationId := loggerCtx.GetCorrelationId()

		details := map[string]string{}

		if requestId := loggerCtx.GetRequestId(); requestId != nil && *requestId != "" {
			details["AWS Request ID"] = *requestId
		}

		if logStream := loggerCtx.GetLogStream(); logStream != nil && *logStream != "" {
			details["Log Stream"] = *logStream
		}

		payload := adapter.NotificationPayload{
			Level:         adapter.NotificationLevelError,
			Title:         "Forest Completion Processing Failed",
			Message:       fmt.Sprintf("Failed to process forest completion.\n\n**Error:** %s", err.Error()),
			TransactionID: correlationId,
			Details:       details,
		}

		if notifyErr := f.notifier.Notify(ctx, payload); notifyErr != nil {
			logger.Warn(ctx, notifyErr, "Failed to send Discord notification")
		}
	}()
}

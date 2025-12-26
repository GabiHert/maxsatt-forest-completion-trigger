package lambda

import (
	"context"
	"errors"
	"fmt"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type errorHandler struct {
	notifier       adapter.Notifier
	failureHandler adapter.FailureHandler
}

func ErrorHandler(notifier adapter.Notifier, failureHandler adapter.FailureHandler) adapter.ErrorHandler {
	return &errorHandler{
		notifier:       notifier,
		failureHandler: failureHandler,
	}
}

func (e *errorHandler) Handle(ctx context.Context, err error, event any) error {
	loggerCtx := logger.GetContext(ctx)

	var customErr errs.BaseError
	if errors.As(err, &customErr) {
		logger.Error(ctx, customErr, "Error processing request", event)
		e.sendErrorNotification(ctx, customErr, loggerCtx.GetCorrelationId())

		if e.failureHandler != nil {
			if dlqErr := e.failureHandler.Handle(ctx, customErr); dlqErr != nil {
				logger.Warn(ctx, dlqErr, "Failed to send to DLQ")
			}
		}

		return customErr
	}

	customErr = errs.InternalServerError(err, "UNK-00500", "Unknown error")
	logger.Error(ctx, customErr, "Unknown error processing request", event)
	e.sendErrorNotification(ctx, customErr, loggerCtx.GetCorrelationId())

	if e.failureHandler != nil {
		if dlqErr := e.failureHandler.Handle(ctx, customErr); dlqErr != nil {
			logger.Warn(ctx, dlqErr, "Failed to send to DLQ")
		}
	}

	return customErr
}

func (e *errorHandler) sendErrorNotification(ctx context.Context, err errs.BaseError, transactionId string) {
	if e.notifier == nil {
		return
	}

	loggerCtx := logger.GetContext(ctx)

	details := map[string]string{
		"Status Code": fmt.Sprintf("%d", err.StatusCode()),
	}

	if correlationId := loggerCtx.GetCorrelationId(); correlationId != "" {
		details["Correlation ID"] = correlationId
	}

	if requestId := loggerCtx.GetRequestId(); requestId != nil && *requestId != "" {
		details["AWS Request ID"] = *requestId
	}

	if logStream := loggerCtx.GetLogStream(); logStream != nil && *logStream != "" {
		details["Log Stream"] = *logStream
	}

	message := err.Error()
	if description := err.Description(); description != "" {
		message = fmt.Sprintf("%s\n\n**Description:** %s", message, description)
	}

	payload := adapter.NotificationPayload{
		Level:         adapter.NotificationLevelError,
		Title:         "Forest Completion Trigger Error",
		Message:       message,
		TransactionID: transactionId,
		ErrorCode:     err.Code(),
		Details:       details,
	}

	if notifyErr := e.notifier.Notify(ctx, payload); notifyErr != nil {
		logger.Warn(ctx, notifyErr, "Failed to send error notification", map[string]any{
			"transactionId": transactionId,
			"errorCode":     err.Code(),
		})
	}
}

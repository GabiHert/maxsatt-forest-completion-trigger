package lambda

import (
	"context"
	"errors"
	"fmt"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/dto"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type errorHandler struct {
	notifier adapter.Notifier
}

// ErrorHandler creates a new error handler middleware.
// The notifier parameter is optional - if nil, no notifications will be sent.
func ErrorHandler(notifier adapter.Notifier) adapter.ErrorHandler {
	return &errorHandler{
		notifier: notifier,
	}
}

func (e *errorHandler) Handle(ctx context.Context, err error, event any) error {
	loggerCtx := logger.GetContext(ctx)
	receiveCount := loggerCtx.GetReceiveCount()
	if receiveCount == 0 {
		receiveCount = 1
	}

	// Check if it's already a BaseError from domain (using pkg/errs constructors)
	var customErr errs.BaseError
	if errors.As(err, &customErr) {
		// Domain errors already use the correct pkg/errs types with proper codes
		logger.Error(ctx, customErr, "Error processing request", event)

		// Send Discord notification
		e.sendErrorNotification(ctx, customErr, loggerCtx.GetCorrelationId(), receiveCount)

		// Handle retries and notifications
		if receiveCount >= customErr.Retries() {
			if customErr.InternalNotify() {
				// Future: additional internal notifications can be added here
			}

			if customErr.Abort() {
				// No abort
			}
		}

		return customErr
	}

	// Unknown error - wrap as internal server error
	customErr = errs.InternalServerError(err, "UNK-00500", "Unknown error")
	logger.Error(ctx, customErr, "Unknown error processing request", event)

	// Send Discord notification for unknown errors
	e.sendErrorNotification(ctx, customErr, loggerCtx.GetCorrelationId(), receiveCount)

	return customErr
}

// sendErrorNotification sends an error notification to the configured notifier.
// If no notifier is configured or if notification fails, it logs a warning but does not fail.
func (e *errorHandler) sendErrorNotification(ctx context.Context, err errs.BaseError, transactionId string, receiveCount int) {
	if e.notifier == nil {
		return
	}

	loggerCtx := logger.GetContext(ctx)

	statusCode := err.StatusCode()
	retriable := statusCode >= 500
	var retryStatus string
	if retriable {
		retryStatus = fmt.Sprintf("Retriable (5xx) - Lambda will retry. Current attempt: %d", receiveCount)
	} else {
		retryStatus = fmt.Sprintf("Non-Retriable (4xx) - Message deleted. Receive count: %d", receiveCount)
	}

	details := map[string]string{
		"Status Code":   fmt.Sprintf("%d", statusCode),
		"Retry Status":  retryStatus,
		"Receive Count": fmt.Sprintf("%d", receiveCount),
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

	if requestValue := ctx.Value("request"); requestValue != nil {
		if req, ok := requestValue.(dto.ClimateEvent); ok {
			details["Processing ID"] = req.ProcessingID
		}
	}

	message := err.Error()
	if description := err.Description(); description != "" {
		message = fmt.Sprintf("%s\n\n**Description:** %s", message, description)
	}

	payload := adapter.NotificationPayload{
		Level:         adapter.NotificationLevelError,
		Title:         "Climate Trigger Error",
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

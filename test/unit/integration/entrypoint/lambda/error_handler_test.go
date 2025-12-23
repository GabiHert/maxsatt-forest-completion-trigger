package lambda_test

import (
	"context"
	"errors"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/dto"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/lambda"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type mockNotifier struct {
	notifyFunc func(ctx context.Context, payload adapter.NotificationPayload) error
	called     bool
	payload    adapter.NotificationPayload
}

func (m *mockNotifier) Notify(ctx context.Context, payload adapter.NotificationPayload) error {
	m.called = true
	m.payload = payload
	if m.notifyFunc != nil {
		return m.notifyFunc(ctx, payload)
	}
	return nil
}

func TestErrorHandler_CustomError(t *testing.T) {
	handler := lambda.ErrorHandler(nil)
	ctx := logger.GetContext(context.Background())

	customErr := errs.BadRequestError("test error", "TEST-001")
	err := handler.Handle(ctx, customErr, nil)

	if err == nil {
		t.Error("Handle() should return error")
	}
}

func TestErrorHandler_NonCustomError(t *testing.T) {
	handler := lambda.ErrorHandler(nil)
	ctx := logger.GetContext(context.Background())

	regularErr := errors.New("regular error")
	err := handler.Handle(ctx, regularErr, nil)

	if err == nil {
		t.Error("Handle() should return error")
	}
}

func TestErrorHandler_WithHighReceiveCount(t *testing.T) {
	handler := lambda.ErrorHandler(nil)

	loggerCtx := logger.GetContext(context.Background())
	loggerCtx.SetReceiveCount(5)

	customErr := errs.BadRequestError("test error", "TEST-001")
	err := handler.Handle(loggerCtx, customErr, map[string]any{"test": "data"})

	if err == nil {
		t.Error("Handle() should return error")
	}
}

func TestErrorHandler_WithNotifier_CustomError(t *testing.T) {
	notifier := &mockNotifier{}
	handler := lambda.ErrorHandler(notifier)

	loggerCtx := logger.GetContext(context.Background())
	loggerCtx.SetReceiveCount(1)

	customErr := errs.BadRequestError("test error", "TEST-001")
	err := handler.Handle(loggerCtx, customErr, nil)

	if err == nil {
		t.Error("Handle() should return error")
	}

	if !notifier.called {
		t.Error("Notifier should have been called")
	}

	if notifier.payload.ErrorCode != "TEST-001" {
		t.Errorf("Expected error code TEST-001, got %s", notifier.payload.ErrorCode)
	}

	if notifier.payload.Level != adapter.NotificationLevelError {
		t.Errorf("Expected notification level error, got %s", notifier.payload.Level)
	}
}

func TestErrorHandler_WithNotifier_NonCustomError(t *testing.T) {
	notifier := &mockNotifier{}
	handler := lambda.ErrorHandler(notifier)

	loggerCtx := logger.GetContext(context.Background())

	regularErr := errors.New("regular error")
	err := handler.Handle(loggerCtx, regularErr, nil)

	if err == nil {
		t.Error("Handle() should return error")
	}

	if !notifier.called {
		t.Error("Notifier should have been called")
	}

	if notifier.payload.ErrorCode != "UNK-00500" {
		t.Errorf("Expected error code UNK-00500, got %s", notifier.payload.ErrorCode)
	}
}

func TestErrorHandler_WithNotifier_ServerError_Retriable(t *testing.T) {
	notifier := &mockNotifier{}
	handler := lambda.ErrorHandler(notifier)

	loggerCtx := logger.GetContext(context.Background())
	loggerCtx.SetReceiveCount(2)

	serverErr := errs.InternalServerError(errors.New("db error"), "DB-001", "Database connection failed")
	err := handler.Handle(loggerCtx, serverErr, nil)

	if err == nil {
		t.Error("Handle() should return error")
	}

	if !notifier.called {
		t.Error("Notifier should have been called")
	}

	if notifier.payload.Details["Status Code"] != "500" {
		t.Errorf("Expected status code 500, got %s", notifier.payload.Details["Status Code"])
	}

	retryStatus := notifier.payload.Details["Retry Status"]
	if retryStatus == "" {
		t.Error("Retry status should be present")
	}
}

func TestErrorHandler_WithNotifier_ClientError_NonRetriable(t *testing.T) {
	notifier := &mockNotifier{}
	handler := lambda.ErrorHandler(notifier)

	loggerCtx := logger.GetContext(context.Background())
	loggerCtx.SetReceiveCount(1)

	clientErr := errs.BadRequestError("invalid input", "VAL-001")
	err := handler.Handle(loggerCtx, clientErr, nil)

	if err == nil {
		t.Error("Handle() should return error")
	}

	if !notifier.called {
		t.Error("Notifier should have been called")
	}

	if notifier.payload.Details["Status Code"] != "400" {
		t.Errorf("Expected status code 400, got %s", notifier.payload.Details["Status Code"])
	}
}

func TestErrorHandler_WithNotifier_NotifyError_DoesNotFail(t *testing.T) {
	notifier := &mockNotifier{
		notifyFunc: func(ctx context.Context, payload adapter.NotificationPayload) error {
			return errors.New("notification failed")
		},
	}
	handler := lambda.ErrorHandler(notifier)

	loggerCtx := logger.GetContext(context.Background())

	customErr := errs.BadRequestError("test error", "TEST-001")
	err := handler.Handle(loggerCtx, customErr, nil)

	if err == nil {
		t.Error("Handle() should return error")
	}

	if !notifier.called {
		t.Error("Notifier should have been called")
	}
}

func TestErrorHandler_WithRequestContext(t *testing.T) {
	notifier := &mockNotifier{}
	handler := lambda.ErrorHandler(notifier)

	loggerCtx := logger.GetContext(context.Background())
	ctxWithRequest := context.WithValue(loggerCtx, "request", dto.ClimateEvent{
		ProcessingID: "proc-123",
	})

	customErr := errs.BadRequestError("test error", "TEST-001")
	err := handler.Handle(ctxWithRequest, customErr, nil)

	if err == nil {
		t.Error("Handle() should return error")
	}

	if notifier.payload.Details["Processing ID"] != "proc-123" {
		t.Errorf("Expected processing ID proc-123, got %s", notifier.payload.Details["Processing ID"])
	}
}

func TestErrorHandler_WithCorrelationId(t *testing.T) {
	notifier := &mockNotifier{}
	handler := lambda.ErrorHandler(notifier)

	loggerCtx := logger.GetContext(context.Background())
	loggerCtx.SetCorrelationId("corr-123")

	customErr := errs.BadRequestError("test error", "TEST-001")
	_ = handler.Handle(loggerCtx, customErr, nil)

	if notifier.payload.Details["Correlation ID"] != "corr-123" {
		t.Errorf("Expected correlation ID corr-123, got %s", notifier.payload.Details["Correlation ID"])
	}

	if notifier.payload.TransactionID != "corr-123" {
		t.Errorf("Expected transaction ID corr-123, got %s", notifier.payload.TransactionID)
	}
}

func TestErrorHandler_WithEmptyRequestIdAndLogStream(t *testing.T) {
	notifier := &mockNotifier{}
	handler := lambda.ErrorHandler(notifier)

	loggerCtx := logger.GetContext(context.Background())

	customErr := errs.BadRequestError("test error", "TEST-001")
	_ = handler.Handle(loggerCtx, customErr, nil)

	// When request ID and log stream are not set, they should not be in details
	if _, ok := notifier.payload.Details["AWS Request ID"]; ok {
		t.Error("AWS Request ID should not be present when not set")
	}

	if _, ok := notifier.payload.Details["Log Stream"]; ok {
		t.Error("Log Stream should not be present when not set")
	}
}

func TestErrorHandler_ZeroReceiveCount_DefaultsToOne(t *testing.T) {
	notifier := &mockNotifier{}
	handler := lambda.ErrorHandler(notifier)

	loggerCtx := logger.GetContext(context.Background())

	customErr := errs.BadRequestError("test error", "TEST-001")
	_ = handler.Handle(loggerCtx, customErr, nil)

	if notifier.payload.Details["Receive Count"] != "1" {
		t.Errorf("Expected receive count 1, got %s", notifier.payload.Details["Receive Count"])
	}
}

func TestErrorHandler_Title_IsCorrect(t *testing.T) {
	notifier := &mockNotifier{}
	handler := lambda.ErrorHandler(notifier)

	loggerCtx := logger.GetContext(context.Background())

	customErr := errs.BadRequestError("test error", "TEST-001")
	_ = handler.Handle(loggerCtx, customErr, nil)

	if notifier.payload.Title != "Climate Trigger Error" {
		t.Errorf("Expected title 'Climate Trigger Error', got %s", notifier.payload.Title)
	}
}

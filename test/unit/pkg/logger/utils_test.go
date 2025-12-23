package logger_test

import (
	"context"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

func TestGetContext_WithBackgroundContext(t *testing.T) {
	ctx := logger.GetContext(context.Background())
	if ctx == nil {
		t.Error("GetContext() should return non-nil context")
	}
}

func TestGetContext_WithNil(t *testing.T) {
	ctx := logger.GetContext(nil)
	if ctx == nil {
		t.Error("GetContext() should return non-nil context even with nil input")
	}
}

func TestGetContext_WithExistingLoggerContext(t *testing.T) {
	ctx1 := logger.GetContext(context.Background())
	ctx2 := logger.GetContext(ctx1)

	if ctx1 != ctx2 {
		t.Error("GetContext() should return same context when passed a logger context")
	}
}

func TestGetContext_GeneratesUniqueIds(t *testing.T) {
	ctx1 := logger.GetContext(context.Background())
	ctx2 := logger.GetContext(context.Background())

	correlationId1 := ctx1.GetCorrelationId()
	correlationId2 := ctx2.GetCorrelationId()

	if correlationId1 == correlationId2 {
		t.Error("GetContext() should generate unique correlation IDs")
	}
}

func TestContext_CorrelationId(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	originalId := ctx.GetCorrelationId()
	if originalId == "" {
		t.Error("GetCorrelationId() should return non-empty string")
	}

	newId := "new-correlation-id"
	ctx.SetCorrelationId(newId)

	if ctx.GetCorrelationId() != newId {
		t.Errorf("SetCorrelationId() failed, expected %s, got %s", newId, ctx.GetCorrelationId())
	}
}

func TestContext_ReceiveCount(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	initialCount := ctx.GetReceiveCount()
	if initialCount != 0 {
		t.Errorf("Initial receive count should be 0, got %d", initialCount)
	}

	ctx.SetReceiveCount(5)
	if ctx.GetReceiveCount() != 5 {
		t.Errorf("SetReceiveCount() failed, expected 5, got %d", ctx.GetReceiveCount())
	}

	ctx.SetReceiveCount(10)
	if ctx.GetReceiveCount() != 10 {
		t.Errorf("SetReceiveCount() should update, expected 10, got %d", ctx.GetReceiveCount())
	}
}

func TestContext_OperationUniqueId(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	initialId := ctx.GetOperationUniqueId()
	if initialId != nil {
		t.Error("Initial operation unique ID should be nil")
	}

	// SetOperationUniqueId only works if operationUniqueId is already set and not empty
	// Since initially it's nil, setting it won't work
	operationId := "operation-123"
	ctx.SetOperationUniqueId(&operationId)

	retrievedId := ctx.GetOperationUniqueId()
	// Should still be nil because the conditional check failed
	if retrievedId != nil {
		t.Error("GetOperationUniqueId() should still be nil because initial value was nil")
	}
}

func TestContext_RequestId(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	requestId := ctx.GetRequestId()
	if requestId != nil {
		t.Logf("Request ID: %s", *requestId)
	}
}

func TestContext_LogGroup(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	logGroup := ctx.GetLogGroup()
	if logGroup != nil {
		t.Logf("Log Group: %s", *logGroup)
	}
}

func TestContext_LogStream(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	logStream := ctx.GetLogStream()
	if logStream != nil {
		t.Logf("Log Stream: %s", *logStream)
	}
}

func TestContext_WithContextValue(t *testing.T) {
	bgCtx := context.Background()
	bgCtx = context.WithValue(bgCtx, "requestId", "test-request-123")
	bgCtx = context.WithValue(bgCtx, "correlationId", "test-correlation-456")
	bgCtx = context.WithValue(bgCtx, "transactionId", "test-transaction-789")

	ctx := logger.GetContext(bgCtx)

	if ctx.GetRequestId() == nil || *ctx.GetRequestId() != "test-request-123" {
		t.Error("Context should preserve requestId from background context")
	}

	if ctx.GetCorrelationId() != "test-correlation-456" {
		t.Error("Context should preserve correlationId from background context")
	}
}

func TestContext_StandardContextMethods(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	deadline, ok := ctx.Deadline()
	if ok {
		t.Logf("Context has deadline: %v", deadline)
	}

	done := ctx.Done()
	if done != nil {
		t.Log("Context has done channel")
	}

	err := ctx.Err()
	if err != nil {
		t.Errorf("Fresh context should not have error, got: %v", err)
	}

	value := ctx.Value("any-key")
	if value != nil {
		t.Logf("Context value: %v", value)
	}
}

func TestContext_Stop(t *testing.T) {
	ctx := logger.GetContext(context.Background())
	ctx.Stop()
}

func TestContext_MultipleOperations(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	ctx.SetCorrelationId("correlation-1")
	ctx.SetReceiveCount(3)
	opId := "op-123"
	ctx.SetOperationUniqueId(&opId)

	if ctx.GetCorrelationId() != "correlation-1" {
		t.Error("Correlation ID should be correlation-1")
	}

	if ctx.GetReceiveCount() != 3 {
		t.Error("Receive count should be 3")
	}

	// SetOperationUniqueId won't work because initial value is nil
	if ctx.GetOperationUniqueId() != nil {
		t.Error("Operation unique ID should still be nil")
	}

	ctx.SetCorrelationId("correlation-2")
	if ctx.GetCorrelationId() != "correlation-2" {
		t.Error("Correlation ID should update to correlation-2")
	}
}

func TestContext_NilOperationUniqueId(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	ctx.SetOperationUniqueId(nil)
	if ctx.GetOperationUniqueId() != nil {
		t.Error("Operation unique ID should be nil after setting to nil")
	}
}

func TestContext_LargeReceiveCount(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	ctx.SetReceiveCount(1000000)
	if ctx.GetReceiveCount() != 1000000 {
		t.Error("Should handle large receive counts")
	}
}

func TestContext_NegativeReceiveCount(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	ctx.SetReceiveCount(-5)
	if ctx.GetReceiveCount() != -5 {
		t.Error("Should handle negative receive counts")
	}
}

func TestContext_EmptyOperationUniqueId(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	emptyId := ""
	ctx.SetOperationUniqueId(&emptyId)

	retrievedId := ctx.GetOperationUniqueId()
	// Should still be nil because initial value was nil
	if retrievedId != nil {
		t.Error("GetOperationUniqueId() should still be nil")
	}
}

func TestContext_ContextInterface(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	underlyingCtx := ctx.Context()
	if underlyingCtx == nil {
		t.Error("Context() should return non-nil context")
	}
}

func TestGetContext_WithMultipleCarriers(t *testing.T) {
	bgCtx := context.Background()
	ctx1 := logger.GetContext(bgCtx)
	ctx2 := logger.GetContext(ctx1)
	ctx3 := logger.GetContext(ctx2)

	if ctx1 == nil || ctx2 == nil || ctx3 == nil {
		t.Error("All contexts should be non-nil")
	}
}

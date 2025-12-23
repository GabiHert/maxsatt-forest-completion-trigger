package logger_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

func TestError_WithMessage(t *testing.T) {
	os.Setenv("LOG_LEVEL", "ERROR")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("test error")

	logger.Error(ctx, err, "Error occurred")
}

func TestError_WithMetadata(t *testing.T) {
	os.Setenv("LOG_LEVEL", "ERROR")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("test error")
	metadata := map[string]string{"key": "value"}

	logger.Error(ctx, err, "Error with metadata", metadata)
}

func TestWarn_WithMessage(t *testing.T) {
	os.Setenv("LOG_LEVEL", "WARN")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("warning error")

	logger.Warn(ctx, err, "Warning occurred")
}

func TestWarn_WithMetadata(t *testing.T) {
	os.Setenv("LOG_LEVEL", "WARN")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("warning error")
	metadata := map[string]int{"count": 5}

	logger.Warn(ctx, err, "Warning with metadata", metadata)
}

func TestInfo_WithMessage(t *testing.T) {
	os.Setenv("LOG_LEVEL", "INFO")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())

	logger.Info(ctx, "Info message")
}

func TestInfo_WithMetadata(t *testing.T) {
	os.Setenv("LOG_LEVEL", "INFO")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	metadata := map[string]bool{"success": true}

	logger.Info(ctx, "Info with metadata", metadata)
}

func TestDebug_WithMessage(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())

	logger.Debug(ctx, "Debug message")
}

func TestDebug_WithMetadata(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	metadata := []string{"item1", "item2"}

	logger.Debug(ctx, "Debug with metadata", metadata)
}

func TestTrace_WithMessage(t *testing.T) {
	os.Setenv("LOG_LEVEL", "TRACE")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())

	logger.Trace(ctx, "Trace message")
}

func TestTrace_WithMetadata(t *testing.T) {
	os.Setenv("LOG_LEVEL", "TRACE")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	metadata := 12345

	logger.Trace(ctx, "Trace with metadata", metadata)
}

func TestLogLevel_Error(t *testing.T) {
	os.Setenv("LOG_LEVEL", "ERROR")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("error")

	logger.Error(ctx, err, "Should print")
	logger.Warn(ctx, err, "Should not print")
	logger.Info(ctx, "Should not print")
	logger.Debug(ctx, "Should not print")
	logger.Trace(ctx, "Should not print")
}

func TestLogLevel_Warn(t *testing.T) {
	os.Setenv("LOG_LEVEL", "WARN")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("error")

	logger.Error(ctx, err, "Should print")
	logger.Warn(ctx, err, "Should print")
	logger.Info(ctx, "Should not print")
	logger.Debug(ctx, "Should not print")
	logger.Trace(ctx, "Should not print")
}

func TestLogLevel_Info(t *testing.T) {
	os.Setenv("LOG_LEVEL", "INFO")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("error")

	logger.Error(ctx, err, "Should print")
	logger.Warn(ctx, err, "Should print")
	logger.Info(ctx, "Should print")
	logger.Debug(ctx, "Should not print")
	logger.Trace(ctx, "Should not print")
}

func TestLogLevel_Debug(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("error")

	logger.Error(ctx, err, "Should print")
	logger.Warn(ctx, err, "Should print")
	logger.Info(ctx, "Should print")
	logger.Debug(ctx, "Should print")
	logger.Trace(ctx, "Should not print")
}

func TestLogLevel_Trace(t *testing.T) {
	os.Setenv("LOG_LEVEL", "TRACE")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("error")

	logger.Error(ctx, err, "Should print")
	logger.Warn(ctx, err, "Should print")
	logger.Info(ctx, "Should print")
	logger.Debug(ctx, "Should print")
	logger.Trace(ctx, "Should print")
}

func TestLogLevel_Off(t *testing.T) {
	os.Setenv("LOG_LEVEL", "OFF")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("error")

	logger.Error(ctx, err, "Should not print")
	logger.Warn(ctx, err, "Should not print")
	logger.Info(ctx, "Should not print")
	logger.Debug(ctx, "Should not print")
	logger.Trace(ctx, "Should not print")
}

func TestLogLevel_Default(t *testing.T) {
	os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("error")

	logger.Error(ctx, err, "Should print with default level")
	logger.Warn(ctx, err, "Should print with default level")
	logger.Info(ctx, "Should print with default level")
	logger.Debug(ctx, "Should print with default level")
	logger.Trace(ctx, "Should print with default level")
}

func TestLogLevel_WithSpaces(t *testing.T) {
	os.Setenv("LOG_LEVEL", " INFO ")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())

	logger.Info(ctx, "Should handle spaces in log level")
}

func TestLogLevel_LowercaseValue(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())

	logger.Info(ctx, "Should handle lowercase log level")
}

func TestError_WithNilError(t *testing.T) {
	os.Setenv("LOG_LEVEL", "ERROR")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())

	logger.Error(ctx, nil, "Error with nil error")
}

func TestError_WithMultipleMetadata(t *testing.T) {
	os.Setenv("LOG_LEVEL", "ERROR")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("test error")
	metadata1 := map[string]string{"key1": "value1"}
	metadata2 := map[string]int{"key2": 42}
	metadata3 := "string metadata"

	logger.Error(ctx, err, "Error with multiple metadata", metadata1, metadata2, metadata3)
}

func TestInfo_WithEmptyMetadata(t *testing.T) {
	os.Setenv("LOG_LEVEL", "INFO")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())

	logger.Info(ctx, "Info with empty metadata")
}

func TestDatadogAgent_Enabled(t *testing.T) {
	os.Setenv("DD_TRACE_ENABLED", "true")
	os.Setenv("DD_AGENT_HOST", "localhost")
	os.Setenv("ENV", "test")
	defer os.Unsetenv("DD_TRACE_ENABLED")
	defer os.Unsetenv("DD_AGENT_HOST")
	defer os.Unsetenv("ENV")

	logger.StartDatadogAgent()
}

func TestDatadogAgent_Disabled(t *testing.T) {
	os.Setenv("DD_TRACE_ENABLED", "false")
	defer os.Unsetenv("DD_TRACE_ENABLED")

	logger.StartDatadogAgent()
}

func TestDatadogAgent_NotSet(t *testing.T) {
	os.Unsetenv("DD_TRACE_ENABLED")

	logger.StartDatadogAgent()
}

func TestLogger_WithDifferentContextTypes(t *testing.T) {
	os.Setenv("LOG_LEVEL", "INFO")
	defer os.Unsetenv("LOG_LEVEL")

	tests := []struct {
		name string
		ctx  any
	}{
		{"with background context", context.Background()},
		{"with logger context", logger.GetContext(context.Background())},
		{"with nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := logger.GetContext(tt.ctx)
			logger.Info(ctx, "Test message")
		})
	}
}

func TestLogger_ContextWithRequestId(t *testing.T) {
	os.Setenv("LOG_LEVEL", "INFO")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := context.WithValue(context.Background(), "requestId", "test-request-id")
	logCtx := logger.GetContext(ctx)

	logger.Info(logCtx, "Message with request ID")

	requestId := logCtx.GetRequestId()
	if requestId == nil || *requestId != "test-request-id" {
		t.Errorf("Expected request ID to be test-request-id, got: %v", requestId)
	}
}

func TestLogger_ContextWithCorrelationId(t *testing.T) {
	os.Setenv("LOG_LEVEL", "INFO")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := context.WithValue(context.Background(), "correlationId", "test-correlation-id")
	logCtx := logger.GetContext(ctx)

	logger.Info(logCtx, "Message with correlation ID")

	correlationId := logCtx.GetCorrelationId()
	if correlationId != "test-correlation-id" {
		t.Errorf("Expected correlation ID to be test-correlation-id, got: %s", correlationId)
	}
}

func TestLogger_ContextWithTransactionId(t *testing.T) {
	os.Setenv("LOG_LEVEL", "INFO")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := context.WithValue(context.Background(), "transactionId", "test-transaction-id")
	logCtx := logger.GetContext(ctx)

	logger.Info(logCtx, "Message with transaction ID")
}

func TestLogger_ContextMethods(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	ctx.SetCorrelationId("new-correlation-id")
	if ctx.GetCorrelationId() != "new-correlation-id" {
		t.Error("SetCorrelationId/GetCorrelationId failed")
	}

	ctx.SetReceiveCount(5)
	if ctx.GetReceiveCount() != 5 {
		t.Error("SetReceiveCount/GetReceiveCount failed")
	}

	operationId := "operation-123"
	ctx.SetOperationUniqueId(&operationId)
	if ctx.GetOperationUniqueId() != nil && *ctx.GetOperationUniqueId() != operationId {
		t.Error("SetOperationUniqueId/GetOperationUniqueId failed")
	}

	if ctx.GetRequestId() != nil {
		t.Log("Request ID:", *ctx.GetRequestId())
	}

	if ctx.GetLogGroup() != nil {
		t.Log("Log Group:", *ctx.GetLogGroup())
	}

	if ctx.GetLogStream() != nil {
		t.Log("Log Stream:", *ctx.GetLogStream())
	}
}

func TestLogger_ContextStandardMethods(t *testing.T) {
	ctx := logger.GetContext(context.Background())

	if _, ok := ctx.Deadline(); ok {
		t.Log("Context has deadline")
	}

	if ctx.Done() != nil {
		t.Log("Context has done channel")
	}

	if ctx.Err() != nil {
		t.Error("Context should not have error initially")
	}

	ctx.Value("test-key")
}

func TestLogger_ContextStop(t *testing.T) {
	ctx := logger.GetContext(context.Background())
	ctx.Stop()
}

func TestLogger_WithLambdaContext(t *testing.T) {
	os.Setenv("AWS_LAMBDA_LOG_GROUP_NAME", "/aws/lambda/test-function")
	os.Setenv("AWS_LAMBDA_LOG_STREAM_NAME", "2025/10/15/[$LATEST]abcd1234")
	defer os.Unsetenv("AWS_LAMBDA_LOG_GROUP_NAME")
	defer os.Unsetenv("AWS_LAMBDA_LOG_STREAM_NAME")

	ctx := logger.GetContext(context.Background())

	logger.Info(ctx, "Lambda context test")
}

func TestError_ComplexMetadata(t *testing.T) {
	os.Setenv("LOG_LEVEL", "ERROR")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("complex error")

	type ComplexStruct struct {
		Field1 string
		Field2 int
		Field3 []string
	}

	metadata := ComplexStruct{
		Field1: "value",
		Field2: 42,
		Field3: []string{"a", "b", "c"},
	}

	logger.Error(ctx, err, "Error with complex metadata", metadata)
}

func TestLogLevel_CaseInsensitive(t *testing.T) {
	tests := []struct {
		level    string
		expected bool
	}{
		{"ERROR", true},
		{"error", true},
		{"Error", true},
		{"WARN", true},
		{"warn", true},
		{"INFO", true},
		{"info", true},
		{"DEBUG", true},
		{"debug", true},
		{"TRACE", true},
		{"trace", true},
		{"OFF", true},
		{"off", true},
		{"invalid", true}, // defaults to true
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", tt.level)
			defer os.Unsetenv("LOG_LEVEL")

			ctx := logger.GetContext(context.Background())
			logger.Info(ctx, "Testing "+tt.level)
		})
	}
}

func TestLogger_ErrorWithDetailedMessage(t *testing.T) {
	os.Setenv("LOG_LEVEL", "ERROR")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("detailed error: something went wrong")

	logger.Error(ctx, err, "Error occurred during processing")
}

func TestLogger_MultipleLogsInSequence(t *testing.T) {
	os.Setenv("LOG_LEVEL", "TRACE")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	err := errors.New("sequence error")

	logger.Trace(ctx, "Starting operation")
	logger.Debug(ctx, "Processing step 1")
	logger.Info(ctx, "Step 1 completed")
	logger.Warn(ctx, err, "Warning in step 2")
	logger.Error(ctx, err, "Error in step 3")
}

func TestLogger_LongMessage(t *testing.T) {
	os.Setenv("LOG_LEVEL", "INFO")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())
	longMessage := strings.Repeat("This is a very long message. ", 100)

	logger.Info(ctx, longMessage)
}

func TestLogger_SpecialCharactersInMessage(t *testing.T) {
	os.Setenv("LOG_LEVEL", "INFO")
	defer os.Unsetenv("LOG_LEVEL")

	ctx := logger.GetContext(context.Background())

	logger.Info(ctx, "Message with special chars: \n\t\"quotes\" and 'apostrophes' and \\ backslashes")
}

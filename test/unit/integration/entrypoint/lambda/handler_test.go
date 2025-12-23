package lambda_test

import (
	"context"
	"errors"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/lambda"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
)

// Mock validator that always passes validation
type mockValidator struct{}

func (m *mockValidator) Struct(_ string, _ interface{}) error {
	return nil
}

type mockErrorHandler struct {
	called bool
}

func (m *mockErrorHandler) Handle(_ context.Context, err error, _ any) error {
	m.called = true
	return err
}

// Mock process climate analysis service that panics
type mockProcessForestCreatedServicePanic struct {
	panicType string
}

func (m *mockProcessForestCreatedServicePanic) Execute(_ context.Context, _ string) (*entity.ClimateAnalysisResult, error) {
	switch m.panicType {
	case "string":
		panic("test panic string")
	case "runtime":
		var x []int
		_ = x[999]
		return nil, nil
	case "error":
		panic(errors.New("test panic error"))
	default:
		panic(123)
	}
}

type mockProcessForestCreatedServiceError struct {
	err error
}

func (m *mockProcessForestCreatedServiceError) Execute(_ context.Context, _ string) (*entity.ClimateAnalysisResult, error) {
	return nil, m.err
}

type mockProcessForestCreatedServiceSuccess struct{}

func (m *mockProcessForestCreatedServiceSuccess) Execute(_ context.Context, _ string) (*entity.ClimateAnalysisResult, error) {
	return &entity.ClimateAnalysisResult{S3Key: "data/processing/test-id/enriched_data.parquet"}, nil
}

func createSQSEvent() events.SQSEvent {
	return events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:              "test-message-id",
				ReceiptHandle:          "test-receipt-handle",
				Body:                   `{"processing_id":"123e4567-e89b-12d3-a456-426614174000"}`,
				Md5OfBody:              "test-md5",
				Md5OfMessageAttributes: "test-md5-attr",
				EventSourceARN:         "arn:aws:sqs:us-east-1:123456789012:test-queue",
				EventSource:            "aws:sqs",
				AWSRegion:              "us-east-1",
				Attributes:             map[string]string{"ApproximateReceiveCount": "1"},
			},
		},
	}
}

func TestHandler_PanicRecovery_String(t *testing.T) {
	errorHandler := &mockErrorHandler{}
	service := &mockProcessForestCreatedServicePanic{panicType: "string"}
	validator := &mockValidator{}
	handler := lambda.Handler(errorHandler, service, validator)

	ctx := logger.GetContext(context.Background())
	event := createSQSEvent()

	_, err := handler.Handle(ctx, event)
	if err == nil {
		t.Error("Handle() should return error when service panics with string")
	}

	if !errorHandler.called {
		t.Error("Error handler should be called when panic occurs")
	}
}

func TestHandler_PanicRecovery_Error(t *testing.T) {
	errorHandler := &mockErrorHandler{}
	service := &mockProcessForestCreatedServicePanic{panicType: "error"}
	validator := &mockValidator{}
	handler := lambda.Handler(errorHandler, service, validator)

	ctx := logger.GetContext(context.Background())
	event := createSQSEvent()

	_, err := handler.Handle(ctx, event)
	if err == nil {
		t.Error("Handle() should return error when service panics with error")
	}

	if !errorHandler.called {
		t.Error("Error handler should be called when panic occurs")
	}
}

func TestHandler_PanicRecovery_Default(t *testing.T) {
	errorHandler := &mockErrorHandler{}
	service := &mockProcessForestCreatedServicePanic{panicType: "default"}
	validator := &mockValidator{}
	handler := lambda.Handler(errorHandler, service, validator)

	ctx := logger.GetContext(context.Background())
	event := createSQSEvent()

	_, err := handler.Handle(ctx, event)
	if err == nil {
		t.Error("Handle() should return error when service panics with default type")
	}

	if !errorHandler.called {
		t.Error("Error handler should be called when panic occurs")
	}
}

func TestHandler_ErrorFromService(t *testing.T) {
	expectedErr := errors.New("service error")
	errorHandler := &mockErrorHandler{}
	service := &mockProcessForestCreatedServiceError{err: expectedErr}
	validator := &mockValidator{}
	handler := lambda.Handler(errorHandler, service, validator)

	ctx := logger.GetContext(context.Background())
	event := createSQSEvent()

	_, err := handler.Handle(ctx, event)
	if err == nil {
		t.Error("Handle() should return error when service returns error")
	}

	if !errorHandler.called {
		t.Error("Error handler should be called when error occurs")
	}
}

func TestHandler_Success(t *testing.T) {
	errorHandler := &mockErrorHandler{}
	service := &mockProcessForestCreatedServiceSuccess{}
	validator := &mockValidator{}
	handler := lambda.Handler(errorHandler, service, validator)

	ctx := logger.GetContext(context.Background())
	event := createSQSEvent()

	_, err := handler.Handle(ctx, event)
	if err != nil {
		t.Errorf("Handle() should not return error on success, got: %v", err)
	}

	if errorHandler.called {
		t.Error("Error handler should not be called on success")
	}
}

func TestHandler_InvalidEvent(t *testing.T) {
	errorHandler := &mockErrorHandler{}
	service := &mockProcessForestCreatedServiceSuccess{}
	validator := &mockValidator{}
	handler := lambda.Handler(errorHandler, service, validator)

	ctx := logger.GetContext(context.Background())
	invalidEvent := "invalid event format"

	_, err := handler.Handle(ctx, invalidEvent)
	if err == nil {
		t.Error("Handle() should return error for invalid event format")
	}
}

func TestHandler_InvalidJSON(t *testing.T) {
	errorHandler := &mockErrorHandler{}
	service := &mockProcessForestCreatedServiceSuccess{}
	validator := &mockValidator{}
	handler := lambda.Handler(errorHandler, service, validator)

	ctx := logger.GetContext(context.Background())
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:      "test-message-id",
				Body:           `{"invalid json"}`,
				EventSourceARN: "arn:aws:sqs:us-east-1:123456789012:test-queue",
				Attributes:     map[string]string{"ApproximateReceiveCount": "1"},
			},
		},
	}

	_, err := handler.Handle(ctx, event)
	if err == nil {
		t.Error("Handle() should return error for invalid JSON")
	}

	if !errorHandler.called {
		t.Error("Error handler should be called for invalid JSON")
	}
}

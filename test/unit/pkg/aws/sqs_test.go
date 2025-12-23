package aws_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type mockSqs struct {
	sendMessageFunc func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

func (m *mockSqs) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	return m.sendMessageFunc(ctx, params, optFns...)
}

type mockLogger struct {
	debugCalled bool
}

func (m *mockLogger) Debug(ctx context.Context, message string, metadata ...any) {
	m.debugCalled = true
}

func (m *mockLogger) GetTransactionID(ctx context.Context) string {
	return "test-transaction-id"
}

func TestSqsClient_WithAWSURL(t *testing.T) {
	os.Setenv("AWS_URL", "http://localhost:4566")
	defer os.Unsetenv("AWS_URL")

	client := aws.SqsClient("us-east-1")
	if client == nil {
		t.Error("SqsClient() should return non-nil client")
	}
}

func TestSqsClient_WithoutAWSURL(t *testing.T) {
	os.Unsetenv("AWS_URL")

	client := aws.SqsClient("us-east-1")
	if client == nil {
		t.Error("SqsClient() should return non-nil client")
	}
}

func TestSqsHelper_Send_Success(t *testing.T) {
	mockSqs := &mockSqs{
		sendMessageFunc: func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			messageId := "test-message-id"
			return &sqs.SendMessageOutput{
				MessageId: &messageId,
			}, nil
		},
	}

	mockLogger := &mockLogger{}
	helper := aws.SqsHelper(mockSqs, mockLogger)

	ctx := context.Background()
	message := map[string]string{"test": "data"}
	queueUrl := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"

	err := helper.Send(ctx, message, queueUrl)
	if err != nil {
		t.Errorf("Send() should not return error, got: %v", err)
	}

	if !mockLogger.debugCalled {
		t.Error("Logger Debug should be called")
	}
}

func TestSqsHelper_Send_MarshalError(t *testing.T) {
	mockSqs := &mockSqs{}
	mockLogger := &mockLogger{}
	helper := aws.SqsHelper(mockSqs, mockLogger)

	ctx := context.Background()
	// Create an unmarshalable value (channels can't be marshaled)
	message := make(chan int)
	queueUrl := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"

	err := helper.Send(ctx, message, queueUrl)
	if err == nil {
		t.Error("Send() should return error for unmarshalable message")
	}
}

func TestSqsHelper_Send_SendMessageError(t *testing.T) {
	expectedError := errors.New("send message failed")
	mockSqs := &mockSqs{
		sendMessageFunc: func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			return nil, expectedError
		},
	}

	mockLogger := &mockLogger{}
	helper := aws.SqsHelper(mockSqs, mockLogger)

	ctx := context.Background()
	message := map[string]string{"test": "data"}
	queueUrl := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"

	err := helper.Send(ctx, message, queueUrl)
	if err == nil {
		t.Error("Send() should return error when SendMessage fails")
	}

	if err != expectedError {
		t.Errorf("Send() should return expected error, got: %v", err)
	}
}

func TestSqsHelper_Send_WithNilMessage(t *testing.T) {
	mockSqs := &mockSqs{
		sendMessageFunc: func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			messageId := "test-message-id"
			return &sqs.SendMessageOutput{
				MessageId: &messageId,
			}, nil
		},
	}

	mockLogger := &mockLogger{}
	helper := aws.SqsHelper(mockSqs, mockLogger)

	ctx := context.Background()
	queueUrl := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"

	err := helper.Send(ctx, nil, queueUrl)
	if err != nil {
		t.Errorf("Send() should handle nil message, got: %v", err)
	}
}

func TestSqsHelper_Send_WithEmptyQueueUrl(t *testing.T) {
	mockSqs := &mockSqs{
		sendMessageFunc: func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			messageId := "test-message-id"
			return &sqs.SendMessageOutput{
				MessageId: &messageId,
			}, nil
		},
	}

	mockLogger := &mockLogger{}
	helper := aws.SqsHelper(mockSqs, mockLogger)

	ctx := context.Background()
	message := map[string]string{"test": "data"}

	err := helper.Send(ctx, message, "")
	if err != nil {
		t.Errorf("Send() should not return error for empty queue URL, got: %v", err)
	}
}

func TestSqsHelper_Send_WithComplexMessage(t *testing.T) {
	mockSqs := &mockSqs{
		sendMessageFunc: func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			messageId := "test-message-id"
			return &sqs.SendMessageOutput{
				MessageId: &messageId,
			}, nil
		},
	}

	mockLogger := &mockLogger{}
	helper := aws.SqsHelper(mockSqs, mockLogger)

	ctx := context.Background()
	message := map[string]any{
		"string":  "value",
		"number":  123,
		"boolean": true,
		"nested": map[string]string{
			"key": "value",
		},
		"array": []string{"item1", "item2"},
	}
	queueUrl := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"

	err := helper.Send(ctx, message, queueUrl)
	if err != nil {
		t.Errorf("Send() should handle complex message, got: %v", err)
	}
}

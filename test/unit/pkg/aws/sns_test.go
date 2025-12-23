package aws_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type mockSns struct {
	publishFunc func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

func (m *mockSns) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	return m.publishFunc(ctx, params, optFns...)
}

func TestSnsClient_WithAWSURL(t *testing.T) {
	os.Setenv("AWS_URL", "http://localhost:4566")
	defer os.Unsetenv("AWS_URL")

	client := aws.SnsClient("us-east-1")
	if client == nil {
		t.Error("SnsClient() should return non-nil client")
	}
}

func TestSnsClient_WithoutAWSURL(t *testing.T) {
	os.Unsetenv("AWS_URL")

	client := aws.SnsClient("us-east-1")
	if client == nil {
		t.Error("SnsClient() should return non-nil client")
	}
}

func TestSnsHelper_Send_Success(t *testing.T) {
	mockSns := &mockSns{
		publishFunc: func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
			messageId := "test-message-id"
			return &sns.PublishOutput{
				MessageId: &messageId,
			}, nil
		},
	}

	mockLogger := &mockLogger{}
	helper := aws.SnsHelper(mockSns, mockLogger)

	ctx := context.Background()
	message := map[string]string{"test": "data"}
	topicArn := "arn:aws:sns:us-east-1:123456789012:test-topic"

	err := helper.Send(ctx, message, topicArn)
	if err != nil {
		t.Errorf("Send() should not return error, got: %v", err)
	}

	if !mockLogger.debugCalled {
		t.Error("Logger Debug should be called")
	}
}

func TestSnsHelper_Send_MarshalError(t *testing.T) {
	mockSns := &mockSns{}
	mockLogger := &mockLogger{}
	helper := aws.SnsHelper(mockSns, mockLogger)

	ctx := context.Background()
	// Create an unmarshalable value (channels can't be marshaled)
	message := make(chan int)
	topicArn := "arn:aws:sns:us-east-1:123456789012:test-topic"

	err := helper.Send(ctx, message, topicArn)
	if err == nil {
		t.Error("Send() should return error for unmarshalable message")
	}
}

func TestSnsHelper_Send_PublishError(t *testing.T) {
	expectedError := errors.New("publish failed")
	mockSns := &mockSns{
		publishFunc: func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
			return nil, expectedError
		},
	}

	mockLogger := &mockLogger{}
	helper := aws.SnsHelper(mockSns, mockLogger)

	ctx := context.Background()
	message := map[string]string{"test": "data"}
	topicArn := "arn:aws:sns:us-east-1:123456789012:test-topic"

	err := helper.Send(ctx, message, topicArn)
	if err == nil {
		t.Error("Send() should return error when Publish fails")
	}

	if err != expectedError {
		t.Errorf("Send() should return expected error, got: %v", err)
	}
}

func TestSnsHelper_Send_WithNilMessage(t *testing.T) {
	mockSns := &mockSns{
		publishFunc: func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
			messageId := "test-message-id"
			return &sns.PublishOutput{
				MessageId: &messageId,
			}, nil
		},
	}

	mockLogger := &mockLogger{}
	helper := aws.SnsHelper(mockSns, mockLogger)

	ctx := context.Background()
	topicArn := "arn:aws:sns:us-east-1:123456789012:test-topic"

	err := helper.Send(ctx, nil, topicArn)
	if err != nil {
		t.Errorf("Send() should handle nil message, got: %v", err)
	}
}

func TestSnsHelper_Send_WithEmptyTopicArn(t *testing.T) {
	mockSns := &mockSns{
		publishFunc: func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
			messageId := "test-message-id"
			return &sns.PublishOutput{
				MessageId: &messageId,
			}, nil
		},
	}

	mockLogger := &mockLogger{}
	helper := aws.SnsHelper(mockSns, mockLogger)

	ctx := context.Background()
	message := map[string]string{"test": "data"}

	err := helper.Send(ctx, message, "")
	if err != nil {
		t.Errorf("Send() should not return error for empty topic ARN, got: %v", err)
	}
}

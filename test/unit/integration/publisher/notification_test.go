package publisher_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/publisher"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type mockSnsHelper struct {
	shouldFail bool
	err        error
}

func (m *mockSnsHelper) Send(ctx context.Context, message any, topic string) error {
	if m.shouldFail {
		return m.err
	}
	return nil
}

func TestNotificationPublisher_PublishNotification_Success(t *testing.T) {
	_ = os.Setenv("FOREST_EVENTS_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:forest-events-topic")
	defer func() {
		_ = os.Unsetenv("FOREST_EVENTS_TOPIC_ARN")
	}()

	mockSns := &mockSnsHelper{shouldFail: false}
	notificationPublisher := publisher.NewNotificationPublisher(mockSns)

	ctx := logger.GetContext(context.Background())
	completion := entity.ForestCompletion{
		ForestId:      "forest-123",
		ProcessingIds: []string{"proc-1", "proc-2", "proc-3"},
	}

	err := notificationPublisher.PublishNotification(ctx, completion)
	if err != nil {
		t.Errorf("PublishNotification() should not return error, got: %v", err)
	}
}

func TestNotificationPublisher_PublishNotification_SnsError(t *testing.T) {
	_ = os.Setenv("FOREST_EVENTS_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:forest-events-topic")
	defer func() {
		_ = os.Unsetenv("FOREST_EVENTS_TOPIC_ARN")
	}()

	expectedErr := errors.New("SNS publish failed")
	mockSns := &mockSnsHelper{shouldFail: true, err: expectedErr}
	notificationPublisher := publisher.NewNotificationPublisher(mockSns)

	ctx := logger.GetContext(context.Background())
	completion := entity.ForestCompletion{
		ForestId:      "forest-123",
		ProcessingIds: []string{"proc-1"},
	}

	err := notificationPublisher.PublishNotification(ctx, completion)
	if err == nil {
		t.Error("PublishNotification() should return error when SNS fails")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("PublishNotification() error = %v, want %v", err, expectedErr)
	}
}

func TestNotificationPublisher_PublishNotification_EmptyProcessingIds(t *testing.T) {
	_ = os.Setenv("FOREST_EVENTS_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:forest-events-topic")
	defer func() {
		_ = os.Unsetenv("FOREST_EVENTS_TOPIC_ARN")
	}()

	mockSns := &mockSnsHelper{shouldFail: false}
	notificationPublisher := publisher.NewNotificationPublisher(mockSns)

	ctx := logger.GetContext(context.Background())
	completion := entity.ForestCompletion{
		ForestId:      "forest-empty",
		ProcessingIds: []string{},
	}

	err := notificationPublisher.PublishNotification(ctx, completion)
	if err != nil {
		t.Errorf("PublishNotification() should not return error for empty processing IDs, got: %v", err)
	}
}

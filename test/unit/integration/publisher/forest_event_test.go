package publisher_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums/eventtype"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/publisher"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type mockSnsHelper struct {
	shouldFail bool
	err        error
}

func (m *mockSnsHelper) Send(_ context.Context, _ any, _ string) error {
	if m.shouldFail {
		return m.err
	}
	return nil
}

func TestForestEventPublisher_Success(t *testing.T) {
	_ = os.Setenv("FOREST_EVENTS_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:test-topic")
	defer func() {
		_ = os.Unsetenv("FOREST_EVENTS_TOPIC_ARN")
	}()

	mockSns := &mockSnsHelper{shouldFail: false}
	publisher := publisher.NewForestEventPublisher(mockSns)

	ctx := logger.GetContext(context.Background())
	err := publisher.Publish(ctx, eventtype.StartAnalysis, map[string]any{"forest_id": "test-id"})
	if err != nil {
		t.Errorf("Publish() should not return error, got: %v", err)
	}
}

func TestForestEventPublisher_SnsError(t *testing.T) {
	_ = os.Setenv("FOREST_EVENTS_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:test-topic")
	defer func() {
		_ = os.Unsetenv("FOREST_EVENTS_TOPIC_ARN")
	}()

	expectedErr := errors.New("SNS publish failed")
	mockSns := &mockSnsHelper{shouldFail: true, err: expectedErr}
	publisher := publisher.NewForestEventPublisher(mockSns)

	ctx := logger.GetContext(context.Background())
	err := publisher.Publish(ctx, eventtype.StartAnalysis, map[string]any{"forest_id": "test-id"})
	if err == nil {
		t.Error("Publish() should return error when SNS fails")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Publish() error = %v, want %v", err, expectedErr)
	}
}

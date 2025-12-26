package publisher

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums/eventtype"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/publisher/dto"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

const (
	notificationPublisherSource = "MAXSATT_FOREST_COMPLETION_TRIGGER"
)

type notificationPublisher struct {
	snsHelper aws.SnsHelperAdapter
}

// NewNotificationPublisher creates a new NotificationPublisher that publishes
// notification events to the forest-events-topic SNS topic.
func NewNotificationPublisher(snsHelper aws.SnsHelperAdapter) adapter.NotificationPublisher {
	return &notificationPublisher{
		snsHelper: snsHelper,
	}
}

// PublishNotification publishes a notification event for a completed forest.
// The event uses the NOTIFY event type and contains the forest_id and processing_ids
// in the event_data field. The forest_id is used as the event_correlation_id.
func (n *notificationPublisher) PublishNotification(ctx context.Context, completion entity.ForestCompletion) error {
	logger.Info(ctx, "Started", map[string]any{
		"forest_id":      completion.ForestId,
		"processing_ids": completion.ProcessingIds,
	})

	eventData := entity.NewNotificationEventData(&completion)

	// Use forest_id as the correlation ID for this notification event
	publishEvent := dto.NewEventWithSource(
		completion.ForestId,
		notificationPublisherSource,
		eventtype.Notify,
		eventData,
	)

	err := n.snsHelper.Send(ctx, publishEvent, properties.Properties().Services.ForestEventsTopicArn)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Finished", publishEvent)
	return nil
}

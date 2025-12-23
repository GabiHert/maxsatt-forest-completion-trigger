package publisher

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/publisher/dto"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type forestEventPublisher struct {
	snsHelper aws.SnsHelperAdapter
}

func NewForestEventPublisher(snsHelper aws.SnsHelperAdapter) adapter.EventPublisher {
	return &forestEventPublisher{
		snsHelper: snsHelper,
	}
}

func (f *forestEventPublisher) Publish(ctx context.Context, eventType enums.EventType, eventData any) error {
	logger.Info(ctx, "Started", map[string]any{
		"event_type": eventType,
		"event_data": eventData,
	})

	correlationId := logger.GetContext(ctx).GetCorrelationId()
	publishEvent := dto.NewEvent(correlationId, eventType, eventData)
	err := f.snsHelper.Send(ctx, publishEvent, properties.Properties().Services.ForestEventsTopicArn)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Finished", publishEvent)
	return nil
}

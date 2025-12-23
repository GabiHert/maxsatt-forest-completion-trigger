package usecase

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums"
)

// PublishEvent defines the use case for publishing events to SNS.
type PublishEvent interface {
	Publish(ctx context.Context, eventType enums.EventType, eventData any) error
}

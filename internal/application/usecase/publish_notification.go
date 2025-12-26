package usecase

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// PublishNotification defines the use case for publishing notification events
// to the forest-events-topic when forests have completed all processings.
type PublishNotification interface {
	// PublishNotification publishes a notification event for a completed forest.
	// The event contains the forest_id and all processing_ids that have completed.
	PublishNotification(ctx context.Context, completion entity.ForestCompletion) error
}

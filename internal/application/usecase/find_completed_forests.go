package usecase

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// FindCompletedForests defines the use case for finding forests that have all
// processings completed and are ready for notification.
type FindCompletedForests interface {
	// FindReadyForNotification finds forests where all processings have status COMPLETED
	// and at least one processing has not been notified (notified_at IS NULL).
	// Returns an empty slice if no forests are ready for notification.
	FindReadyForNotification(ctx context.Context, limit int) ([]entity.ForestCompletion, error)
}

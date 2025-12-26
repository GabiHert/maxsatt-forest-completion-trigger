package usecase

import (
	"context"
)

// MarkAsNotified defines the use case for marking processings as notified
// after a successful notification has been published to the forest-events-topic.
type MarkAsNotified interface {
	// MarkAsNotified updates the notified_at timestamp for the given processing IDs.
	// This operation is idempotent - only processings with notified_at IS NULL will be updated.
	// Returns an error if the database operation fails.
	MarkAsNotified(ctx context.Context, processingIds []string) error
}

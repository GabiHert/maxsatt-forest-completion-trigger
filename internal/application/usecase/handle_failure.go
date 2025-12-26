package usecase

import (
	"context"
)

// HandleFailure defines the use case for handling processing failures.
// It sends failed records to a dead letter queue and notifies relevant systems.
type HandleFailure interface {
	HandleFailure(ctx context.Context, forestId string, processingIds []string, err error) error
}

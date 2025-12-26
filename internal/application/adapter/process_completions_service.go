package adapter

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// ProcessCompletionsService defines the interface for the main forest completion processing service.
// This service orchestrates the full flow of finding completed forests, publishing notifications,
// and marking processings as notified.
type ProcessCompletionsService interface {
	// Execute processes forest completions by:
	// 1. Finding forests where all processings are completed and not yet notified
	// 2. For each forest: publishing a notification event to SNS
	// 3. On success: marking all processings for that forest as notified
	// Returns a result indicating how many forests were processed and any errors encountered.
	Execute(ctx context.Context, maxForests int) (*entity.CompletionResult, error)
}

package adapter

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"
)

// ForestCompletionRepository defines the interface for querying forests ready for notification
// and marking processings as notified after successful notification.
type ForestCompletionRepository interface {
	usecase.FindCompletedForests
	usecase.MarkAsNotified
}

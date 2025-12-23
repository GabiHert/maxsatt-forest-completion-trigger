package adapter

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"
)

// EventPublisher defines the interface for publishing forest events.
type EventPublisher interface {
	usecase.PublishEvent
}

package adapter

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"
)

// NotificationPublisher defines the interface for publishing notification events
// to the forest-events-topic SNS topic.
type NotificationPublisher interface {
	usecase.PublishNotification
}

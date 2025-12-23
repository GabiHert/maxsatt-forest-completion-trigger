package adapter

import "context"

type NotificationLevel string

const (
	NotificationLevelError   NotificationLevel = "error"
	NotificationLevelWarning NotificationLevel = "warning"
	NotificationLevelInfo    NotificationLevel = "info"
)

type NotificationPayload struct {
	Level         NotificationLevel
	Title         string
	Message       string
	TransactionID string
	ErrorCode     string
	Details       map[string]string
}

// Notifier defines the interface for sending notifications to external systems.
// Implementations should log errors gracefully without failing the application.
type Notifier interface {
	Notify(ctx context.Context, payload NotificationPayload) error
}

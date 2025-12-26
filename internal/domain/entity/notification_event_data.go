package entity

// NotificationEventData represents the event_data payload for the notification event
// that is published to the forest-events-topic SNS topic.
// This data structure is used in the event_data field of the CloudEvent message.
type NotificationEventData struct {
	// ForestId is the unique identifier of the forest that has all processings completed.
	ForestId string `json:"forest_id"`

	// ProcessingIds contains all processing IDs that belong to this forest
	// and have been marked as completed.
	ProcessingIds []string `json:"processing_ids"`
}

// NewNotificationEventData creates a new NotificationEventData from a ForestCompletion entity.
func NewNotificationEventData(forestCompletion *ForestCompletion) *NotificationEventData {
	return &NotificationEventData{
		ForestId:      forestCompletion.ForestId,
		ProcessingIds: forestCompletion.ProcessingIds,
	}
}

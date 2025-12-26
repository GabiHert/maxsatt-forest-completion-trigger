package entity

// ForestCompletion represents a forest that has all processings completed
// and is ready for notification. It contains the forest ID and the list
// of processing IDs that belong to this forest.
type ForestCompletion struct {
	// ForestId is the unique identifier of the forest.
	ForestId string

	// ProcessingIds contains all processing IDs that belong to this forest
	// and have status = COMPLETED with notified_at IS NULL.
	ProcessingIds []string
}

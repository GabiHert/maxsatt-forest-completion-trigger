package entity

import "time"

// ProcessingMetadata represents the metadata associated with a processing job.
// This metadata is derived from the Analysis API response and parquet delta data,
// NOT from a separate S3 metadata.json file.
type ProcessingMetadata struct {
	Date         time.Time
	CreatedAt    time.Time
	ProcessingID string
	FieldName    string
	FieldID      string
	Status       string
}

// NewProcessingMetadata creates a new ProcessingMetadata from Analysis API data and date range.
func NewProcessingMetadata(analysisData *Analysis) *ProcessingMetadata {
	return &ProcessingMetadata{
		ProcessingID: analysisData.ProcessingID,
		FieldName:    analysisData.FieldName,
		FieldID:      analysisData.FieldID,
		Date:         analysisData.Date,
		Status:       analysisData.Status,
		CreatedAt:    time.Now(),
	}
}

// Package dto contains data transfer objects for the entrypoint layer.
package dto

// ClimateEvent represents the event payload for climate analysis requests.
type ClimateEvent struct {
	ProcessingID string `json:"processing_id" validate:"required,uuid4"`
}

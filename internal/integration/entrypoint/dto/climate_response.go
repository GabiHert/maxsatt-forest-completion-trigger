// Package dto contains data transfer objects for the entrypoint layer.
package dto

// ClimateResponse represents the response payload for climate analysis.
type ClimateResponse struct {

	// Skipped indicates if processing was skipped due to date validation.
	// When true, S3Key will be empty and Reason will contain the skip reason.
	Skipped *bool `json:"skipped,omitempty"`

	Metadata *ResponseMetadata `json:"metadata,omitempty"`
	// S3Key contains the S3 key where the enriched parquet file was written.
	// The enriched dataset preserves all original parquet fields plus added weather columns.
	// This streaming approach avoids loading the entire dataset into memory.
	S3Key string `json:"s3Key,omitempty"`

	// Reason contains the reason for skipping processing when Skipped is true.
	// Possible values: "no image data available", "climate analysis is up to date"
	Reason string `json:"reason,omitempty"`
}

// ResponseMetadata contains execution metadata about the climate analysis.
type ResponseMetadata struct {
	CacheHit         *bool  `json:"cacheHit,omitempty"`
	Message          string `json:"message"`
	ExecutionTimeMs  int64  `json:"executionTimeMs"`
	WeatherAPITimeMs int64  `json:"weatherApiTimeMs"`
	MergingTimeMs    int64  `json:"mergingTimeMs"`
}

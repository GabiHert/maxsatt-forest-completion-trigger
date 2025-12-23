package entity

// ClimateAnalysisResult represents the result of a climate analysis execution,
// including the S3 location of the enriched data and execution metrics.
type ClimateAnalysisResult struct {
	// S3Key contains the S3 key where the enriched parquet file was written.
	// The enriched dataset contains all original parquet fields plus weather columns.
	// This streaming approach avoids loading the entire dataset into memory.
	S3Key string

	// Skipped indicates if processing was skipped due to date validation.
	// When true, S3Key will be empty and Reason will contain the skip reason.
	Skipped *bool

	// Reason contains the reason for skipping processing when Skipped is true.
	// Possible values: "no image data available", "climate analysis is up to date"
	Reason string

	// Metadata contains execution metrics
	Metadata ExecutionMetadata
}

// ExecutionMetadata contains timing and performance metrics for the analysis execution.
type ExecutionMetadata struct {
	CacheHit         *bool
	WeatherAPITimeMs int64
	MergingTimeMs    int64
}

// NewSkippedResult creates a ClimateAnalysisResult indicating processing was skipped.
func NewSkippedResult(reason string) *ClimateAnalysisResult {
	skipped := true
	return &ClimateAnalysisResult{
		Skipped: &skipped,
		Reason:  reason,
	}
}

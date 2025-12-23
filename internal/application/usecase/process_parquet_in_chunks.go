package usecase

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// ProcessParquetInChunks defines the use case for streaming parquet data processing.
// This interface allows processing large parquet files without loading all data into memory.
type ProcessParquetInChunks interface {
	// ProcessByProcessingId fetches a parquet file from S3, processes it in chunks,
	// and writes enriched results directly to S3 without accumulating data in memory.
	// Returns the S3 key where the enriched parquet file was written.
	ProcessByProcessingId(ctx context.Context, processingID string, analysisData *entity.Analysis, weatherMetrics *entity.WeatherMetrics) (string, error)
}

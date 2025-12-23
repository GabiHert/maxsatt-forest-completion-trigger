package adapter

import (
	"context"
	"io"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// ParquetProcessor defines the interface for reading parquet files and writing enriched results.
// This abstraction allows the use of different parquet libraries without affecting business logic.
// The processor streams data to avoid loading entire datasets into memory.
type ParquetProcessor interface {
	// ProcessInChunks reads parquet data from a reader, enriches it with weather data,
	// and writes the results incrementally to S3 using the provided writer.
	// This method processes row groups incrementally without accumulating all data in memory.
	// Returns the S3 key where the enriched parquet file was written.
	ProcessInChunks(ctx context.Context, reader io.Reader, analysisData *entity.Analysis, weatherMetrics *entity.WeatherMetrics, writer ParquetWriter) (string, error)
}

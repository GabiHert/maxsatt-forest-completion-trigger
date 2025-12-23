package adapter

import (
	"context"
)

// ParquetWriter defines the interface for writing parquet data incrementally to S3.
// This abstraction allows building a parquet file by appending chunks without accumulating
// all data in memory. The writer handles row group optimization and S3 upload.
type ParquetWriter interface {
	// AppendRows adds rows to the parquet file being built.
	// Rows are buffered until enough accumulate to form a row group, then written.
	// This method is called multiple times during streaming processing.
	AppendRows(ctx context.Context, rows []map[string]any) error

	// Finalize completes the parquet file and uploads it to S3.
	// It writes any remaining buffered rows, finalizes the parquet structure,
	// and returns the S3 key where the file was written.
	// This method should be called once after all AppendRows calls complete.
	Finalize(ctx context.Context) (string, error)
}

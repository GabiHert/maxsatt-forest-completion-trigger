package utils

import (
	"bytes"
	"context"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/parquet-go/parquet-go"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

// parquetWriter implements the ParquetWriter interface for streaming parquet writes to S3.
// It uses parquet-go library for schema-based writing with map[string]any data.
// The writer uses an explicit schema derived from the input parquet file schema,
// extended with additional weather metric columns.
type parquetWriter struct {
	s3Client     aws.S3
	schema       *parquet.Schema
	bucketName   string
	s3Key        string
	buffer       []map[string]any
	allRows      []map[string]any
	rowGroupSize int
}

// NewParquetWriterWithS3Client creates a new ParquetWriter with an explicit S3 client and schema.
// The schema parameter should contain the full schema including original fields + weather metrics.
// rowGroupSize determines how many rows to buffer before writing a row group (recommended: 1000-10000).
func NewParquetWriterWithS3Client(
	s3Client aws.S3,
	bucketName string,
	s3Key string,
	rowGroupSize int,
) adapter.ParquetWriter {
	return &parquetWriter{
		s3Client:     s3Client,
		bucketName:   bucketName,
		s3Key:        s3Key,
		buffer:       make([]map[string]any, 0, rowGroupSize),
		rowGroupSize: rowGroupSize,
		allRows:      []map[string]any{}, // Dynamic growth to handle any dataset size
	}
}

// SetSchema sets the parquet schema for this writer.
// This must be called before AppendRows to establish the schema.
// The schema should be the extended schema (input + weather columns).
func (w *parquetWriter) SetSchema(schema *parquet.Schema) {
	w.schema = schema
}

// AppendRows adds rows to the parquet file being built.
// Rows are accumulated in memory until Finalize is called.
// This method can be called multiple times during streaming processing.
func (w *parquetWriter) AppendRows(ctx context.Context, rows []map[string]any) error {
	logger.Debug(ctx, "Started", "rowCount", len(rows), "currentBufferSize", len(w.buffer))

	w.allRows = append(w.allRows, rows...)

	logger.Debug(ctx, "Finished", "totalRows", len(w.allRows))
	return nil
}

// Finalize completes the parquet file and uploads it to S3.
// It writes all accumulated rows to a parquet file using the established schema,
// then uploads the complete file to S3.
func (w *parquetWriter) Finalize(ctx context.Context) (string, error) {
	logger.Debug(ctx, "Started", "totalRows", len(w.allRows))

	if len(w.allRows) == 0 {
		return "", fmt.Errorf("no rows to write")
	}

	if w.schema == nil {
		return "", fmt.Errorf("schema not set - call SetSchema before writing")
	}

	buffer := new(bytes.Buffer)

	writer := parquet.NewGenericWriter[map[string]any](
		buffer,
		w.schema,
		parquet.PageBufferSize(256*1024),
	)

	_, err := writer.Write(w.allRows)
	if err != nil {
		return "", errs.InternalServerError(err, "PARQUET-WRITER-01-0001", "failed to write rows to parquet")
	}

	if err = writer.Close(); err != nil {
		return "", errs.InternalServerError(err, "PARQUET-WRITER-01-0002", "failed to close parquet writer")
	}

	putInput := &s3.PutObjectInput{
		Bucket: awssdk.String(w.bucketName),
		Key:    awssdk.String(w.s3Key),
		Body:   bytes.NewReader(buffer.Bytes()),
	}

	if _, err = w.s3Client.PutObject(ctx, putInput); err != nil {
		return "", errs.InternalServerError(err, "PARQUET-WRITER-01-0003", "failed to upload parquet to S3")
	}

	logger.Info(ctx, "Finished", "bucket", w.bucketName, "key", w.s3Key, "size", buffer.Len())
	return w.s3Key, nil
}

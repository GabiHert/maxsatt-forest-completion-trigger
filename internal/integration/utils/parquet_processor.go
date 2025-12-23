package utils

import (
	"bytes"
	"context"
	"io"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"

	"github.com/parquet-go/parquet-go"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

// parquetProcessor implements the ParquetProcessor interface using a streaming approach.
// It uses github.com/parquet-go/parquet-go for reading parquet files dynamically.
// This implementation processes data in chunks and writes immediately to avoid memory accumulation.
type parquetProcessor struct {
	mergeDatasets usecase.MergeDatasets
}

// NewParquetProcessor creates a new ParquetProcessor implementation.
func NewParquetProcessor(mergeDatasets usecase.MergeDatasets) adapter.ParquetProcessor {
	return &parquetProcessor{
		mergeDatasets: mergeDatasets,
	}
}

// ProcessInChunks reads parquet data from a reader, enriches it with weather data,
// and writes results incrementally to S3 using the provided writer.
// This method is memory-efficient as it processes row groups incrementally and writes
// immediately without accumulating all rows in memory. Returns the S3 key of the output file.
func (p *parquetProcessor) ProcessInChunks(ctx context.Context, reader io.Reader, analysisData *entity.Analysis, weatherMetrics *entity.WeatherMetrics, writer adapter.ParquetWriter) (string, error) {
	logger.Debug(ctx, "Starting", analysisData, weatherMetrics)

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", errs.InternalServerError(err, "PARQUET-02-0001", "failed to read parquet data")
	}

	bytesReader := bytes.NewReader(data)
	parquetFile, err := parquet.OpenFile(bytesReader, int64(len(data)))
	if err != nil {
		return "", errs.InternalServerError(err, "PARQUET-02-0002", "failed to open parquet file")
	}

	inputSchema := parquetFile.Schema()
	extendedSchema := p.buildExtendedSchema(inputSchema)

	if schemaWriter, ok := writer.(interface{ SetSchema(*parquet.Schema) }); ok {
		schemaWriter.SetSchema(extendedSchema)
	}

	const chunkSize = 500
	rowGroups := parquetFile.RowGroups()
	chunk := make([]map[string]any, 0, chunkSize)
	processedRows := int64(0)
	metadata := entity.NewProcessingMetadata(analysisData)

	chunkProcessor := func(chunk []map[string]any) error {
		mergedChunk, err := p.mergeDatasets.Merge(ctx, chunk, weatherMetrics, metadata)
		if err != nil {
			return err
		}

		if err := writer.AppendRows(ctx, mergedChunk); err != nil {
			return err
		}

		return nil
	}

	for _, rowGroup := range rowGroups {
		numRowsInGroup := rowGroup.NumRows()

		rowGroupReader := parquet.NewRowGroupReader(rowGroup)

		for i := int64(0); i < numRowsInGroup; i++ {
			row := make(map[string]any)
			err = rowGroupReader.Read(&row)
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", errs.InternalServerError(err, "PARQUET-02-0003", "failed to read parquet row")
			}

			chunk = append(chunk, row)
			processedRows++

			if len(chunk) >= chunkSize {
				if err = chunkProcessor(chunk); err != nil {
					return "", errs.InternalServerError(err, "PARQUET-02-0004", "failed to process chunk")
				}
				chunk = make([]map[string]any, 0, chunkSize)
			}
		}
	}

	if len(chunk) > 0 {
		if err = chunkProcessor(chunk); err != nil {
			return "", errs.InternalServerError(err, "PARQUET-02-0005", "failed to process final chunk")
		}
	}

	s3Key, err := writer.Finalize(ctx)
	if err != nil {
		return "", errs.InternalServerError(err, "PARQUET-02-0006", "failed to finalize parquet file")
	}

	logger.Info(ctx, "Finished processing parquet data in chunks", "totalRows", processedRows, "s3Key", s3Key)
	return s3Key, nil
}

// buildExtendedSchema creates a new schema by reconstructing the input schema and adding weather metric fields.
// This allows the output parquet file to contain all original columns plus the new weather data columns.
// Fields must be reconstructed rather than copied to ensure compatibility with GenericWriter.
func (p *parquetProcessor) buildExtendedSchema(inputSchema *parquet.Schema) *parquet.Schema {
	inputFields := inputSchema.Fields()

	group := make(parquet.Group)

	for _, field := range inputFields {
		fieldType := field.Type()
		fieldNode := parquet.Leaf(fieldType)

		if field.Optional() {
			fieldNode = parquet.Optional(fieldNode)
		} else if field.Repeated() {
			fieldNode = parquet.Repeated(fieldNode)
		} else if field.Required() {
			fieldNode = parquet.Required(fieldNode)
		}

		group[field.Name()] = fieldNode
	}

	group["avg_temperature"] = parquet.Optional(parquet.Leaf(parquet.DoubleType))
	group["temp_std_dev"] = parquet.Optional(parquet.Leaf(parquet.DoubleType))
	group["avg_humidity"] = parquet.Optional(parquet.Leaf(parquet.DoubleType))
	group["humidity_std_dev"] = parquet.Optional(parquet.Leaf(parquet.DoubleType))
	group["total_precipitation"] = parquet.Optional(parquet.Leaf(parquet.DoubleType))
	group["dry_days_consecutive"] = parquet.Optional(parquet.Leaf(parquet.Int32Type))

	return parquet.NewSchema("enriched_data", group)
}

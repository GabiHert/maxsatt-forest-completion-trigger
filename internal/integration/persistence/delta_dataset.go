package persistence

import (
	"context"
	"fmt"
	"strings"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	domainError "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/error"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/utils"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

// deltaDataSetPersistence implements the FetchDeltaDataset usecase.
type deltaDataSetPersistence struct {
	s3Client         aws.S3
	s3Helper         aws.S3HelperAdapter
	parquetProcessor adapter.ParquetProcessor
}

// NewDeltaDataset creates a new FetchDeltaDataset usecase implementation.
// s3Client must implement the aws.S3 interface.
func NewDeltaDataset(
	s3Client aws.S3,
	s3Helper aws.S3HelperAdapter,
	parquetProcessor adapter.ParquetProcessor,
) adapter.DeltaDatasetPersistenceAdapter {
	return &deltaDataSetPersistence{
		s3Client:         s3Client,
		s3Helper:         s3Helper,
		parquetProcessor: parquetProcessor,
	}
}

// ProcessByProcessingId fetches delta dataset from S3, processes it in chunks with weather enrichment,
// and writes the enriched data directly to S3 without accumulating in memory.
// This method is memory-efficient as it streams data from S3, processes incrementally,
// and writes results immediately. Returns the S3 key where enriched parquet was written.
// Returns ProcessingNotFoundError if the input parquet file does not exist.
// Path structure: processings/{forest_id}/{field_id}/{date(yyyy-MM-dd)}/
// Input files: tries clean.parquet first, falls back to raw.parquet
// Output file: climate.parquet
func (u *deltaDataSetPersistence) ProcessByProcessingId(ctx context.Context, processingID string, analysisData *entity.Analysis, weatherMetrics *entity.WeatherMetrics) (string, error) {
	bucketName := properties.Properties().Services.ClimateDataBucket
	basePath := analysisData.BasePath()

	s3Object, err := u.fetchInputParquet(ctx, bucketName, basePath, processingID)
	if err != nil {
		return "", err
	}
	defer s3Object.Reader.Close()

	outputParquetKey := fmt.Sprintf("%s/climate.parquet", basePath)

	parquetWriter := utils.NewParquetWriterWithS3Client(u.s3Client, bucketName, outputParquetKey, 5000)

	s3Key, err := u.parquetProcessor.ProcessInChunks(ctx, s3Object.Reader, analysisData, weatherMetrics, parquetWriter)
	if err != nil {
		return "", err
	}

	return s3Key, nil
}

// fetchInputParquet attempts to fetch the input parquet file from S3.
// Returns ProcessingNotFoundError if the file does not exist.
func (u *deltaDataSetPersistence) fetchInputParquet(ctx context.Context, bucketName, basePath, processingID string) (*aws.Object, error) {
	deltaParquetKey := fmt.Sprintf("%s/delta.parquet", basePath)

	s3Object, err := u.s3Helper.GetObject(ctx, bucketName, deltaParquetKey)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") || strings.Contains(err.Error(), "not found") {
			return nil, domainError.NewProcessingNotFoundError(processingID)
		}
		logger.Error(ctx, err, "S3 operation failed while fetching delta.parquet", "bucket", bucketName, "key", deltaParquetKey)
		return nil, errs.InternalServerError(err, "DELTA-02-0001", "failed to fetch delta dataset")
	}

	return s3Object, nil
}

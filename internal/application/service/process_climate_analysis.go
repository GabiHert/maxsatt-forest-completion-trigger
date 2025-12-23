package service

import (
	"context"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums/eventtype"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums/filetype"
	domainError "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/error"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type processClimateAnalysisService struct {
	processParquetInChunks  usecase.ProcessParquetInChunks
	fetchFieldAnalysis      usecase.FetchFieldAnalysis
	fetchWeatherData        usecase.FetchWeatherData
	calculateWeatherMetrics usecase.CalculateWeatherMetrics
	fetchDeltaFile          usecase.FetchDeltaFile
	createFileRecord        usecase.CreateFileRecord
	publishEvent            usecase.PublishEvent
}

func NewProcessClimateAnalysisService(
	processParquetInChunks usecase.ProcessParquetInChunks,
	fetchFieldAnalysis usecase.FetchFieldAnalysis,
	fetchWeatherData usecase.FetchWeatherData,
	calculateWeatherMetrics usecase.CalculateWeatherMetrics,
	fetchDeltaFile usecase.FetchDeltaFile,
	createFileRecord usecase.CreateFileRecord,
	publishEvent usecase.PublishEvent,
) adapter.ProcessClimateAnalysisService {
	return &processClimateAnalysisService{
		processParquetInChunks:  processParquetInChunks,
		fetchFieldAnalysis:      fetchFieldAnalysis,
		fetchWeatherData:        fetchWeatherData,
		calculateWeatherMetrics: calculateWeatherMetrics,
		fetchDeltaFile:          fetchDeltaFile,
		createFileRecord:        createFileRecord,
		publishEvent:            publishEvent,
	}
}

// Execute processes a climate analysis for the given processing ID.
// Input validation is performed at the entry point (handler layer) before reaching this service.
// Writes enriched parquet file directly to S3 without loading dataset into memory.
// Uses streaming approach to process parquet data in chunks and write results incrementally.
// Returns the S3 key where the enriched parquet file was written.
// If date validation determines processing should be skipped, returns a skipped result with nil error.
func (s *processClimateAnalysisService) Execute(ctx context.Context, processingID string) (*entity.ClimateAnalysisResult, error) {
	logger.Info(ctx, "Started", processingID)

	analysisData, err := s.fetchFieldAnalysis.FetchByProcessingId(ctx, processingID)
	if err != nil {
		return nil, err
	}

	if shouldSkip, reason := analysisData.ShouldSkipProcessing(); shouldSkip {
		logger.Info(ctx, "Skipping processing", "reason", reason, "processingID", processingID)
		return entity.NewSkippedResult(reason), nil
	}

	if !analysisData.IsCentroidValid() {
		return nil, domainError.NewInvalidCoordinatesError(analysisData.Latitude(), analysisData.Longitude())
	}

	weatherAPITimeStart := time.Now()
	weatherData, err := s.fetchWeatherData.FetchWeather(ctx, analysisData.Latitude(), analysisData.Longitude(), analysisData.Date)
	weatherAPITimeMs := time.Since(weatherAPITimeStart).Milliseconds()
	if err != nil {
		return nil, err
	}

	weatherMetrics, err := s.calculateWeatherMetrics.Calculate(ctx, weatherData, analysisData.Date)
	if err != nil {
		return nil, err
	}

	mergingStart := time.Now()
	s3Key, err := s.processParquetInChunks.ProcessByProcessingId(ctx, processingID, analysisData, weatherMetrics)
	if err != nil {
		return nil, err
	}
	mergingTimeMs := time.Since(mergingStart).Milliseconds()

	if s3Key == "" {
		return nil, domainError.NewEmptyDatasetError()
	}

	deltaFile, err := s.fetchDeltaFile.GetByProcessingID(ctx, processingID)
	if err != nil {
		return nil, err
	}

	fileRecord := entity.NewFileRecord(deltaFile.ImageID, filetype.Climate, s3Key)
	_, err = s.createFileRecord.CreateFile(ctx, fileRecord)
	if err != nil {
		return nil, err
	}

	err = s.publishEvent.Publish(ctx, eventtype.ReadyAnalysis, map[string]string{
		"processing_id": processingID,
	})
	if err != nil {
		return nil, err
	}

	cacheHit := false
	result := &entity.ClimateAnalysisResult{
		S3Key: s3Key,
		Metadata: entity.ExecutionMetadata{
			WeatherAPITimeMs: weatherAPITimeMs,
			MergingTimeMs:    mergingTimeMs,
			CacheHit:         &cacheHit,
		},
	}

	logger.Info(ctx, "Finished", "s3Key", s3Key, "weatherAPITimeMs", weatherAPITimeMs, "mergingTimeMs", mergingTimeMs)
	return result, nil
}

package usecase

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	domainError "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/error"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

// MergeDatasets defines the use case for combining delta data with weather metrics.
type MergeDatasets interface {
	// Merge merges delta pixel data with weather metrics.
	// Works with generic map data to preserve all original fields and add weather columns.
	// Returns enriched data as maps ready for parquet writing.
	// All rows receive the same weather metrics as they are from the same analysis date.
	Merge(ctx context.Context, deltaData []map[string]any, weatherMetrics *entity.WeatherMetrics, metadata *entity.ProcessingMetadata) ([]map[string]any, error)
}

// mergeDatasets implements the MergeDatasets usecase.
type mergeDatasets struct {
}

// NewMergeDatasetsUseCase creates a new MergeDatasets usecase implementation.
func NewMergeDatasetsUseCase() MergeDatasets {
	return &mergeDatasets{}
}

// Merge merges delta pixel data with weather metrics.
// Applies the same weather metrics to all rows since they are from the same analysis date.
// Preserves all original fields from delta data and adds weather metric columns.
func (m *mergeDatasets) Merge(ctx context.Context, deltaData []map[string]any, weatherMetrics *entity.WeatherMetrics, metadata *entity.ProcessingMetadata) ([]map[string]any, error) {
	logger.Info(ctx, "Started", deltaData, weatherMetrics, metadata)

	if metadata == nil {
		return nil, domainError.NewNilMetadataError()
	}

	if weatherMetrics == nil {
		return nil, domainError.NewNilWeatherMetricsError()
	}

	enrichedData := make([]map[string]any, 0, len(deltaData))

	for _, row := range deltaData {
		if row == nil {
			continue
		}

		enrichedRow := make(map[string]any)
		for k, v := range row {
			enrichedRow[k] = v
		}

		enrichedRow["avg_temperature"] = weatherMetrics.AvgTemperature
		enrichedRow["temp_std_dev"] = weatherMetrics.TempStdDev
		enrichedRow["avg_humidity"] = weatherMetrics.AvgHumidity
		enrichedRow["humidity_std_dev"] = weatherMetrics.HumidityStdDev
		enrichedRow["total_precipitation"] = weatherMetrics.TotalPrecipitation
		enrichedRow["dry_days_consecutive"] = weatherMetrics.DryDaysConsecutive

		enrichedData = append(enrichedData, enrichedRow)
	}

	logger.Info(ctx, "Finished", enrichedData)
	return enrichedData, nil
}

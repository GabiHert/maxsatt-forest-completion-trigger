package adapter

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// ProcessClimateAnalysisService defines the interface for the main climate analysis service.
type ProcessClimateAnalysisService interface {
	// Execute processes a climate analysis for the given processing ID.
	// It fetches data from S3, calculates weather metrics, and returns enriched final data with execution metadata.
	Execute(ctx context.Context, processingID string) (*entity.ClimateAnalysisResult, error)
}

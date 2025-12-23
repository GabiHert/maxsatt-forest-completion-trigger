package usecase

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// FetchFieldAnalysis defines the use case for retrieving field analysis data from the API.
type FetchFieldAnalysis interface {
	// FetchByProcessingId fetches analysis data from the Analysis API.
	FetchByProcessingId(ctx context.Context, processingID string) (*entity.Analysis, error)
}

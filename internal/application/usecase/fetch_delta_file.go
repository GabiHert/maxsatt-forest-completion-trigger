package usecase

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// FetchDeltaFile defines the use case for fetching a DELTA file by processing ID.
type FetchDeltaFile interface {
	// GetByProcessingID fetches the DELTA file for a processing to get the image_id.
	// Returns the file record or an error if not found.
	GetByProcessingID(ctx context.Context, processingID string) (*entity.FileRecord, error)
}

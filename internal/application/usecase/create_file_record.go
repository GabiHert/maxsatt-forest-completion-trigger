package usecase

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// CreateFileRecord defines the use case for creating a file record in the external API.
type CreateFileRecord interface {
	// CreateFile creates a file record in the external API to register the created parquet file.
	// It receives the image ID, file type, and file path to register.
	// Returns the created file record entity or an error if creation fails.
	CreateFile(ctx context.Context, file *entity.FileRecord) (*entity.FileRecord, error)
}

package adapter

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"
)

// FileWebService defines the interface for file operations in the external API.
type FileWebService interface {
	usecase.CreateFileRecord
	usecase.FetchDeltaFile
}

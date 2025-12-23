package entity

import "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums"

// FileRecord represents a file record to be registered in the external API.
// It contains the metadata needed to register a parquet file after creation.
type FileRecord struct {
	// ID is the unique identifier of the file record (returned by the API after creation).
	ID string

	// ImageID is the identifier of the image/processing associated with this file.
	ImageID string

	// Type is the type of the file (e.g., Climate, Raw, Clean, Tif).
	Type enums.FileType

	// FilePath is the S3 path where the file is stored.
	FilePath string
}

// NewFileRecord creates a new FileRecord with the required fields.
func NewFileRecord(imageID string, fileType enums.FileType, filePath string) *FileRecord {
	return &FileRecord{
		ImageID:  imageID,
		Type:     fileType,
		FilePath: filePath,
	}
}

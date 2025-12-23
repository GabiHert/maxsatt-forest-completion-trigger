package dto

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums/filetype"
)

// FileRecord represents the file creation API.
type FileRecord struct {
	ID       string `json:"id"`
	ImageID  string `json:"image_id"`
	Type     string `json:"type"`
	FilePath string `json:"file_path"`
}

// NewFileRecord creates a new FileRecord from a domain entity.
func NewFileRecord(file *entity.FileRecord) *FileRecord {
	return &FileRecord{
		ImageID:  file.ImageID,
		Type:     file.Type.String(),
		FilePath: file.FilePath,
	}
}

// ToEntity converts the API response to a domain entity.
func (r *FileRecord) ToEntity() *entity.FileRecord {
	return &entity.FileRecord{
		ID:       r.ID,
		ImageID:  r.ImageID,
		Type:     filetype.GetFileType(r.Type),
		FilePath: r.FilePath,
	}
}

// PaginationMetadata represents pagination information in list responses.
type PaginationMetadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ListFilesResponse represents the response from the files list API.
type ListFilesResponse struct {
	Content  []FileRecord       `json:"content"`
	Files    []FileRecord       `json:"files"`
	Data     []FileRecord       `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

// GetFiles returns the files from the response, supporting multiple API response formats.
func (l *ListFilesResponse) GetFiles() []FileRecord {
	if len(l.Content) > 0 {
		return l.Content
	}
	if len(l.Files) > 0 {
		return l.Files
	}
	return l.Data
}

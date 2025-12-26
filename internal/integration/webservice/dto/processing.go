package dto

import (
	"context"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

// ProcessingResponse represents a processing returned from the MaxSatt API
type ProcessingResponse struct {
	ID         string     `json:"id"`
	FieldID    string     `json:"field_id"`
	ForestID   string     `json:"forest_id,omitempty"`
	Status     string     `json:"status"`
	NotifiedAt *time.Time `json:"notified_at,omitempty"`
	CreatedAt  string     `json:"created_at"`
	UpdatedAt  string     `json:"updated_at"`
	Field      *FieldInfo `json:"field,omitempty"`
}

// FieldInfo represents embedded field information in the processing response
type FieldInfo struct {
	ID       string `json:"id"`
	ForestID string `json:"forest_id"`
	Name     string `json:"name"`
}

// ListProcessingsResponse represents the response from listing processings
type ListProcessingsResponse struct {
	Data       []ProcessingResponse `json:"data"`
	Pagination *PaginationResponse  `json:"pagination,omitempty"`
}

// PaginationResponse represents pagination information
type PaginationResponse struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
}

// UpdateProcessingRequest represents the request to update a processing
type UpdateProcessingRequest struct {
	NotifiedAt *time.Time `json:"notified_at,omitempty"`
	Status     *string    `json:"status,omitempty"`
}

// NewUpdateProcessingNotifiedRequest creates a request to mark a processing as notified
func NewUpdateProcessingNotifiedRequest(notifiedAt time.Time) UpdateProcessingRequest {
	return UpdateProcessingRequest{
		NotifiedAt: &notifiedAt,
	}
}

// ToForestCompletions groups processings by forest and converts to domain entities.
// It logs a warning if any processings are found without a forest_id (orphaned).
func ToForestCompletions(ctx context.Context, processings []ProcessingResponse) []entity.ForestCompletion {
	forestMap := make(map[string][]string)
	var orphanedProcessings []string

	for _, p := range processings {
		forestID := p.ForestID
		if forestID == "" && p.Field != nil {
			forestID = p.Field.ForestID
		}

		if forestID == "" {
			orphanedProcessings = append(orphanedProcessings, p.ID)
			continue
		}

		forestMap[forestID] = append(forestMap[forestID], p.ID)
	}

	if len(orphanedProcessings) > 0 {
		logger.Warn(ctx, nil, "Found processings without forest_id - possible data integrity issue", map[string]any{
			"orphanedProcessingIDs": orphanedProcessings,
			"count":                 len(orphanedProcessings),
		})
	}

	result := make([]entity.ForestCompletion, 0, len(forestMap))
	for forestID, processingIDs := range forestMap {
		result = append(result, entity.ForestCompletion{
			ForestId:      forestID,
			ProcessingIds: processingIDs,
		})
	}

	return result
}

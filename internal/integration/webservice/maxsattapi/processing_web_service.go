package maxsattapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/webservice/dto"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

// ProcessingWebService implements adapter.ForestCompletionRepository using MaxSatt API
type ProcessingWebService struct {
	client *MaxsattAPIClient
}

// NewProcessingWebService creates a new ProcessingWebService instance
func NewProcessingWebService(client *MaxsattAPIClient) adapter.ForestCompletionRepository {
	return &ProcessingWebService{
		client: client,
	}
}

// FindReadyForNotification finds forests where all processings have status COMPLETED
// and at least one processing has not been notified (notified_at IS NULL).
// It calls GET /v1/processings?status=COMPLETED&notified_at_is_null=true&include=field
// and groups the results by forest_id.
func (p *ProcessingWebService) FindReadyForNotification(ctx context.Context, limit int) ([]entity.ForestCompletion, error) {
	logger.Debug(ctx, "Started", map[string]any{
		"limit": limit,
	})

	// Build the query path with parameters
	path := fmt.Sprintf("/v1/processings?status=COMPLETED&notified_at_is_null=true&include=field&per_page=%d", limit)

	var respBody dto.ListProcessingsResponse
	err := p.client.DoRequest(ctx, http.MethodGet, path, nil, &respBody)
	if err != nil {
		logger.Error(ctx, err, "Failed to fetch processings ready for notification via API")
		return nil, err
	}

	if len(respBody.Data) == 0 {
		logger.Debug(ctx, "No processings ready for notification")
		return []entity.ForestCompletion{}, nil
	}

	// Group processings by forest and convert to domain entities
	result := dto.ToForestCompletions(ctx, respBody.Data)

	logger.Debug(ctx, "Finished", map[string]any{
		"processingCount": len(respBody.Data),
		"forestCount":     len(result),
	})

	return result, nil
}

// MarkAsNotified updates the notified_at timestamp for the given processing IDs.
// It calls PUT /v1/processing/{id} with notified_at timestamp for each processing.
// This operation is idempotent - only processings with notified_at IS NULL will be updated.
// Uses fail-fast strategy: returns immediately on first error to preserve at-most-once semantics.
func (p *ProcessingWebService) MarkAsNotified(ctx context.Context, processingIds []string) error {
	logger.Debug(ctx, "Started", map[string]any{
		"processingIds": processingIds,
		"count":         len(processingIds),
	})

	if len(processingIds) == 0 {
		logger.Debug(ctx, "No processing IDs to mark as notified")
		return nil
	}

	notifiedAt := time.Now().UTC()
	reqBody := dto.NewUpdateProcessingNotifiedRequest(notifiedAt)

	// NOTE: This function uses fail-fast error handling intentionally.
	// Partial failure handling was considered but rejected because:
	// 1. The application layer expects all-or-nothing semantics for notification marking
	// 2. On retry, the application re-fetches processings with notified_at IS NULL
	// 3. If we marked some as notified before failing, retry would skip those already marked
	// 4. This would cause "ghost" notifications - SNS messages sent for forests where
	//    MarkAsNotified appeared to fail, but actually partially succeeded
	// The correct behavior is: fail fast, let the entire operation retry from scratch
	for _, processingID := range processingIds {
		path := fmt.Sprintf("/v1/processing/%s", processingID)

		var respBody dto.ProcessingResponse
		err := p.client.DoRequest(ctx, http.MethodPut, path, reqBody, &respBody)
		if err != nil {
			logger.Error(ctx, err, "Failed to mark processing as notified via API", map[string]any{
				"processingID": processingID,
			})
			return err
		}

		logger.Debug(ctx, "Marked processing as notified", map[string]any{
			"processingID": processingID,
			"notifiedAt":   notifiedAt,
		})
	}

	logger.Debug(ctx, "Finished", map[string]any{
		"markedCount": len(processingIds),
	})

	return nil
}

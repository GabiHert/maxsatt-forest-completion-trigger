package service

import (
	"context"
	"fmt"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type processCompletionsService struct {
	findCompletedForests usecase.FindCompletedForests
	publishNotification  usecase.PublishNotification
	markAsNotified       usecase.MarkAsNotified
}

func NewProcessCompletionsService(
	findCompletedForests usecase.FindCompletedForests,
	publishNotification usecase.PublishNotification,
	markAsNotified usecase.MarkAsNotified,
) adapter.ProcessCompletionsService {
	return &processCompletionsService{
		findCompletedForests: findCompletedForests,
		publishNotification:  publishNotification,
		markAsNotified:       markAsNotified,
	}
}

func (s *processCompletionsService) Execute(ctx context.Context, maxForests int) (*entity.CompletionResult, error) {
	logger.Info(ctx, "Started", map[string]any{
		"maxForests": maxForests,
	})

	result := entity.NewCompletionResult()

	forests, err := s.findCompletedForests.FindReadyForNotification(ctx, maxForests)
	if err != nil {
		return nil, fmt.Errorf("failed to find completed forests: %w", err)
	}

	if len(forests) == 0 {
		logger.Info(ctx, "No forests found ready for notification")
		return result, nil
	}

	logger.Info(ctx, "Found forests ready for notification", map[string]any{
		"count": len(forests),
	})

	for _, forest := range forests {
		err := s.processForest(ctx, forest)
		if err != nil {
			logger.Error(ctx, err, "Failed to process forest", map[string]any{
				"forestId": forest.ForestId,
			})
			result.AddError(fmt.Errorf("forest %s: %w", forest.ForestId, err))
			continue
		}

		result.IncrementProcessed()
	}

	logger.Info(ctx, "Finished", map[string]any{
		"processedCount": result.ProcessedCount,
		"failedCount":    result.FailedCount,
		"errorsCount":    len(result.Errors),
	})

	return result, nil
}

func (s *processCompletionsService) processForest(ctx context.Context, forest entity.ForestCompletion) error {
	logger.Debug(ctx, "Processing forest", map[string]any{
		"forestId":        forest.ForestId,
		"processingCount": len(forest.ProcessingIds),
	})

	// Mark as notified FIRST to prevent duplicate notifications on retry.
	// This provides "at-most-once" semantics. If publish fails after marking,
	// the forest won't be re-queried, but the error is logged for monitoring.
	// Trade-off: Prefer no duplicates over missed notifications for this use case.
	err := s.markAsNotified.MarkAsNotified(ctx, forest.ProcessingIds)
	if err != nil {
		return fmt.Errorf("failed to mark as notified: %w", err)
	}

	err = s.publishNotification.PublishNotification(ctx, forest)
	if err != nil {
		logger.Error(ctx, err, "Failed to publish after marking - notification may be lost", map[string]any{
			"forestId":      forest.ForestId,
			"processingIds": forest.ProcessingIds,
		})
		return fmt.Errorf("failed to publish notification: %w", err)
	}

	logger.Debug(ctx, "Forest processed successfully", map[string]any{
		"forestId": forest.ForestId,
	})

	return nil
}

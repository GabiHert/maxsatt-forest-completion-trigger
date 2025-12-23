package lambda

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/adapter"
	integrationAdapter "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/dto"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/validator"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type handler struct {
	errorHandler                  integrationAdapter.ErrorHandler
	processClimateAnalysisService adapter.ProcessClimateAnalysisService
	validator                     validator.Validate
}

func Handler(
	errorHandler integrationAdapter.ErrorHandler,
	processClimateAnalysisService adapter.ProcessClimateAnalysisService,
	validate validator.Validate,
) integrationAdapter.Handler {
	return &handler{
		errorHandler:                  errorHandler,
		processClimateAnalysisService: processClimateAnalysisService,
		validator:                     validate,
	}
}

func (h *handler) Handle(lambdaCtx context.Context, event any) (responseData any, err error) {
	ctx := logger.GetContext(lambdaCtx)
	startTime := time.Now()
	logger.Info(ctx, "Started", event)

	defer func(ctx context.Context, event any) {
		if recovered := recover(); recovered != nil {
			err = deferred(ctx, event, recovered)
		}
		if err != nil {
			err = h.errorHandler.Handle(ctx, err, event)
		}
	}(ctx, event)

	var climateEvent dto.ClimateEvent
	_, retryCount, err := aws.SqsEventParser(event, &climateEvent)
	if err != nil {
		return nil, err
	}

	ctx.SetReceiveCount(retryCount)
	ctx.SetCorrelationId(climateEvent.ProcessingID)

	if err = h.validator.Struct("CLIMATE-01-0001", &climateEvent); err != nil {
		logger.Error(ctx, err, "Event validation failed")
		return nil, err
	}

	result, err := h.processClimateAnalysisService.Execute(ctx, climateEvent.ProcessingID)
	if err != nil {
		return nil, err
	}

	if result.Skipped != nil && *result.Skipped {
		response := dto.ClimateResponse{
			Skipped: result.Skipped,
			Reason:  result.Reason,
		}
		logger.Info(ctx, "Finished - Processing skipped", response)
		return response, nil
	}

	response := dto.ClimateResponse{
		S3Key: result.S3Key,
		Metadata: &dto.ResponseMetadata{
			ExecutionTimeMs:  time.Since(startTime).Milliseconds(),
			WeatherAPITimeMs: result.Metadata.WeatherAPITimeMs,
			MergingTimeMs:    result.Metadata.MergingTimeMs,
			CacheHit:         result.Metadata.CacheHit,
			Message:          "Climate analysis completed successfully",
		},
	}

	logger.Info(ctx, "Finished", response)
	return response, nil
}

func deferred(ctx context.Context, event any, recovered any) error {
	var message string
	switch e := recovered.(type) {
	case string:
		message = fmt.Sprint("recovered (string) panic:", e)
	case runtime.Error:
		message = fmt.Sprint("recovered (runtime.Error) panic:", e.Error())
	case error:
		message = fmt.Sprint("recovered (error) panic:", e.Error())
	default:
		message = fmt.Sprint("recovered (default) panic:", e)
	}

	err := errs.PanicError(message)

	logger.Error(ctx, err, "Panic occurred on lambda execution", event)
	return err
}

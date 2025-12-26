package lambda

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/adapter"
	integrationAdapter "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

const maxForestsPerInvocation = 100

type handler struct {
	errorHandler              integrationAdapter.ErrorHandler
	processCompletionsService adapter.ProcessCompletionsService
}

func Handler(
	errorHandler integrationAdapter.ErrorHandler,
	processCompletionsService adapter.ProcessCompletionsService,
) integrationAdapter.Handler {
	return &handler{
		errorHandler:              errorHandler,
		processCompletionsService: processCompletionsService,
	}
}

func (h *handler) Handle(lambdaCtx context.Context, event any) (responseData any, err error) {
	ctx := logger.GetContext(lambdaCtx)
	startTime := time.Now()
	logger.Info(ctx, "Started forest completion trigger", event)

	defer func(ctx context.Context, event any) {
		if recovered := recover(); recovered != nil {
			err = deferred(ctx, event, recovered)
		}
		if err != nil {
			err = h.errorHandler.Handle(ctx, err, event)
		}
	}(ctx, event)

	result, err := h.processCompletionsService.Execute(ctx, maxForestsPerInvocation)
	if err != nil {
		return nil, err
	}

	response := map[string]any{
		"processed_count":   result.ProcessedCount,
		"failed_count":      result.FailedCount,
		"execution_time_ms": time.Since(startTime).Milliseconds(),
	}

	if result.ProcessedCount == 0 && result.FailedCount == 0 {
		response["message"] = "No forests ready for notification"
	} else {
		response["message"] = fmt.Sprintf("Processed %d forests, %d failed", result.ProcessedCount, result.FailedCount)
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

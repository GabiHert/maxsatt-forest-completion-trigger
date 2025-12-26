package application

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/dependency"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"

	"github.com/aws/aws-lambda-go/lambda"
)

func Start() {
	logger.Info(context.TODO(), "Starting forest completion trigger lambda")

	start := dependency.Injector().Wire(context.Background()).Handler.Handle
	lambda.Start(func(ctx context.Context, event any) error {
		_, err := start(ctx, event)
		return err
	})
}

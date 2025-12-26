package main

import (
	"context"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/dependency"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/google/uuid"
)

func main() {
	awsCtx := lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{
		AwsRequestID: uuid.New().String(),
	})

	// Forest completion trigger is scheduler-based, so the event is empty
	event := map[string]any{}

	response, err := dependency.Injector().Wire(awsCtx).Handler.Handle(awsCtx, event)
	if err != nil {
		panic(err)
	}

	println("Response:", response)
}

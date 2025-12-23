package main

import (
	"context"
	"encoding/json"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/dependency"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/google/uuid"
)

func main() {

	// Sample climate analysis event
	message := map[string]any{
		"processing_id": "test-processing-123",
	}

	awsCtx := lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{
		AwsRequestID: uuid.New().String(),
	})

	messageBodyBytes, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	messageEvent := &events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:              "",
				ReceiptHandle:          "",
				Body:                   string(messageBodyBytes),
				Md5OfBody:              "",
				Md5OfMessageAttributes: "",
				Attributes:             nil,
				MessageAttributes:      nil,
				EventSourceARN:         "",
				EventSource:            "",
				AWSRegion:              "",
			},
		},
	}

	response, err := dependency.Injector().Wire(awsCtx).Handler.Handle(awsCtx, messageEvent)
	if err != nil {
		panic(err)
	}

	println("Response:", response)
}

package aws

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsHelperAdapter interface {
	Send(ctx context.Context, message any, queueUrl string) error
}

type Sqs interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type sqsHelper struct {
	logger loggerAdapter
	sqs    Sqs
}

func SqsClient(region string) *sqs.Client {
	cfg := getConfig(region)
	if awsUrl := os.Getenv("AWS_URL"); awsUrl != "" {
		return sqs.NewFromConfig(cfg, func(o *sqs.Options) {
			o.BaseEndpoint = &awsUrl
		})
	}

	return sqs.NewFromConfig(cfg)
}

func SqsHelper(sqs Sqs, logger loggerAdapter) SqsHelperAdapter {
	return &sqsHelper{
		logger: logger,
		sqs:    sqs,
	}
}

func (s *sqsHelper) Send(ctx context.Context, message any, queueUrl string) error {
	s.logger.Debug(ctx, "Started", message, queueUrl)

	stringMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	payload := string(stringMessage)
	messageInput := &sqs.SendMessageInput{
		MessageBody: &payload,
		QueueUrl:    &queueUrl,
	}

	result, err := s.sqs.SendMessage(ctx, messageInput)
	if err != nil {
		return err
	}

	s.logger.Debug(ctx, "Finished", messageInput, result)
	return nil
}

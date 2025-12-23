package aws

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SnsHelperAdapter interface {
	Send(ctx context.Context, message any, topic string) error
}

type Sns interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type snsHelper struct {
	logger loggerAdapter
	sns    Sns
}

func SnsClient(region string) *sns.Client {
	cfg := getConfig(region)
	if awsUrl := os.Getenv("AWS_URL"); awsUrl != "" {
		return sns.NewFromConfig(cfg, func(o *sns.Options) {
			o.BaseEndpoint = &awsUrl
		})
	}

	return sns.NewFromConfig(cfg)
}

func SnsHelper(sns Sns, logger loggerAdapter) SnsHelperAdapter {
	return &snsHelper{
		logger: logger,
		sns:    sns,
	}
}

func (s *snsHelper) Send(ctx context.Context, message any, topicArn string) error {
	s.logger.Debug(ctx, "Started", message, topicArn)

	stringMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	payload := string(stringMessage)
	messageInput := &sns.PublishInput{
		Message:  &payload,
		TopicArn: &topicArn,
	}

	result, err := s.sns.Publish(ctx, messageInput)
	if err != nil {
		return err
	}

	s.logger.Debug(ctx, "Finished", messageInput, result)
	return nil
}

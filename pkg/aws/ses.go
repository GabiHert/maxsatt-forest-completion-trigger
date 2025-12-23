package aws

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SesHelperAdapter interface {
	SendTemplatedEmail(ctx context.Context, source, toEmail, template, configurationSet string, templateData map[string]string) error
}

type Ses interface {
	SendTemplatedEmail(ctx context.Context, params *ses.SendTemplatedEmailInput, optFns ...func(*ses.Options)) (*ses.SendTemplatedEmailOutput, error)
}

type sesHelper struct {
	logger loggerAdapter
	ses    Ses
	client ses.Client
}

func SesClient(region string) *ses.Client {
	cfg := getConfig(region)
	if awsUrl := os.Getenv("AWS_URL"); awsUrl != "" {
		return ses.NewFromConfig(cfg, func(o *ses.Options) {
			o.BaseEndpoint = &awsUrl
		})
	}

	return ses.NewFromConfig(cfg)
}

func SesHelper(ses Ses, logger loggerAdapter) SesHelperAdapter {
	return &sesHelper{
		logger: logger,
		ses:    ses,
	}
}

func (s *sesHelper) SendTemplatedEmail(ctx context.Context, source, toEmail, template, configurationSet string, templateData map[string]string) error {
	s.logger.Debug(ctx, "Started", map[string]any{
		"source":           source,
		"toEmail":          toEmail,
		"template":         template,
		"configurationSet": configurationSet,
		"templateData":     templateData,
	})

	templateDataBytes, err := json.Marshal(templateData)
	if err != nil {
		return err
	}

	result, err := s.ses.SendTemplatedEmail(ctx, &ses.SendTemplatedEmailInput{
		Source: &source,
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Template:             &template,
		TemplateData:         aws.String(string(templateDataBytes)),
		ConfigurationSetName: &configurationSet,
		Tags: []types.MessageTag{
			{
				Name:  aws.String("EmailType"),
				Value: aws.String("ClientRegistration"),
			},
		},
	})
	if err != nil {
		return err
	}

	s.logger.Debug(ctx, "Finished", map[string]any{
		"source":           source,
		"toEmail":          toEmail,
		"template":         template,
		"configurationSet": configurationSet,
		"templateData":     templateData,
		"result":           result,
	})
	return nil
}

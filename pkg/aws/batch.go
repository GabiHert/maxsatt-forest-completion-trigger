package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/batch/types"
)

type BatchHelperAdapter interface {
	Submit(ctx context.Context, jobName, jobQueue, jobDefinition string, command []string) error
}

type Batch interface {
	SubmitJob(ctx context.Context, params *batch.SubmitJobInput, optFns ...func(*batch.Options)) (*batch.SubmitJobOutput, error)
}

type batchHelper struct {
	logger loggerAdapter
	batch  Batch
}

func BatchClient(region string) *batch.Client {
	cfg := getConfig(region)
	if awsUrl := os.Getenv("AWS_URL"); awsUrl != "" {
		return batch.NewFromConfig(cfg, func(o *batch.Options) {
			o.BaseEndpoint = &awsUrl
		})
	}

	return batch.NewFromConfig(cfg)
}

func BatchHelper(batch Batch, logger loggerAdapter) BatchHelperAdapter {
	return &batchHelper{
		logger: logger,
		batch:  batch,
	}
}

func (b *batchHelper) Submit(ctx context.Context, jobName, jobQueue, jobDefinition string, command []string) error {
	b.logger.Debug(ctx, "Started", map[string]any{
		"jobName":       jobName,
		"jobQueue":      jobQueue,
		"jobDefinition": jobDefinition,
		"command":       command,
	})

	input := &batch.SubmitJobInput{
		JobName:       &jobName,
		JobQueue:      &jobQueue,
		JobDefinition: &jobDefinition,
		ContainerOverrides: &types.ContainerOverrides{
			Command: command,
		},
	}
	result, err := b.batch.SubmitJob(ctx, input)
	if err != nil {
		return err
	}

	b.logger.Debug(ctx, "Finished", input, result)
	return nil
}

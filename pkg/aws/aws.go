package aws

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	aws2 "github.com/lsgndln/dd-trace-go/contrib/aws/aws-sdk-go-v2/aws"
)

type loggerAdapter interface {
	Debug(ctx context.Context, message string, metadata ...any)
	GetTransactionID(ctx context.Context) string
}

type redisAdapter interface {
	Put(ctx context.Context, id string, object any, optionalDuration ...time.Duration) error
	Get(ctx context.Context, id string, object any) error
}

var (
	awsConfig map[string]*aws.Config
)

// getConfig creates aws configuration
func getConfig(region string) aws.Config {
	if awsConfig == nil {
		awsConfig = map[string]*aws.Config{}
	}
	if awsConfig[region] == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			panic(err)
		}
		awsConfig[region] = &cfg
	}

	if trace := os.Getenv("DD_TRACE_ENABLED"); trace == "true" {
		aws2.AppendMiddleware(awsConfig[region])
	}

	return *awsConfig[region]
}

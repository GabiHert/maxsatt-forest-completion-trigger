package dependency

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/service"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	integrationAdapter "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/lambda"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/messaging"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/publisher"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/secret"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/webservice/maxsattapi"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/http"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
	redishelper "github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/redis"
)

type injector struct {
	loggerWrapper logger.Logger

	// AWS Clients
	Sns            aws.Sns
	Sqs            aws.Sqs
	SecretsManager aws.SecretsManager

	// Helpers
	snsHelper             aws.SnsHelperAdapter
	sqsHelper             aws.SqsHelperAdapter
	redisHelper           redishelper.RedisHelper
	secretsManagerHelper  aws.SecretsManagerHelperAdapter
	secretsManagerAdapter integrationAdapter.SecretsManagerAdapter

	// External clients
	httpClient       http.Client
	maxsattAPIClient *maxsattapi.MaxsattAPIClient

	// Publishers
	notificationPublisher integrationAdapter.NotificationPublisher

	// Repositories
	forestCompletionRepository integrationAdapter.ForestCompletionRepository

	// Notifiers
	notifier integrationAdapter.Notifier

	// Failure handling
	failureHandler integrationAdapter.FailureHandler

	// Services
	processCompletionsService adapter.ProcessCompletionsService

	// Entrypoint
	ErrorHandler integrationAdapter.ErrorHandler
	Handler      integrationAdapter.Handler
}

var injectorInit sync.Once
var instance *injector

func Injector() *injector {
	if instance == nil {
		injectorInit.Do(
			func() {
				instance = &injector{}
			},
		)
	}
	return instance
}

// ResetInjector resets the injector singleton - should only be used in tests
func ResetInjector() {
	instance = nil
	injectorInit = sync.Once{}
}

func (i *injector) Wire(ctx context.Context) *injector {
	if i.loggerWrapper == nil {
		i.loggerWrapper = logger.LoggerWrapper()
	}

	// Initialize Secrets Manager first (needed for properties)
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = os.Getenv("AWS_SNS_REGION")
	}

	if i.SecretsManager == nil {
		i.SecretsManager = aws.SecretsManagerClient(awsRegion)
	}

	// Use NoOp Redis Helper since we don't need caching for secrets in this service
	if i.redisHelper == nil {
		i.redisHelper = redishelper.NoOpRedisHelper()
	}

	if i.secretsManagerHelper == nil {
		i.secretsManagerHelper = aws.SecretsManagerHelper(i.SecretsManager, i.redisHelper, i.loggerWrapper)
	}

	if i.secretsManagerAdapter == nil {
		i.secretsManagerAdapter = secret.NewSecretsManager(i.secretsManagerHelper)
	}

	// Initialize properties with secrets manager
	properties.InitializeSecretsManager(i.secretsManagerAdapter)

	props := properties.Properties()

	if i.Sns == nil {
		i.Sns = aws.SnsClient(props.Aws.Sns.Region)
	}

	if i.snsHelper == nil {
		i.snsHelper = aws.SnsHelper(i.Sns, i.loggerWrapper)
	}

	if i.Sqs == nil {
		i.Sqs = aws.SqsClient(props.Aws.Sqs.Region)
	}

	if i.sqsHelper == nil {
		i.sqsHelper = aws.SqsHelper(i.Sqs, i.loggerWrapper)
	}

	if i.httpClient == nil {
		i.httpClient = http.HttpClient(30*time.Second, i.loggerWrapper)
	}

	// Initialize MaxSatt API client
	if i.maxsattAPIClient == nil {
		apiProps := props.MaxsattAPI
		i.maxsattAPIClient = maxsattapi.NewMaxsattAPIClient(
			apiProps.BaseURL,
			apiProps.AuthURL,
			apiProps.ClientID,
			apiProps.ClientSecret,
			i.httpClient,
			i.loggerWrapper,
		)
		logger.Info(ctx, "MaxSatt API client initialized", map[string]any{
			"baseURL": apiProps.BaseURL,
			"authURL": apiProps.AuthURL,
		})
	}

	if i.notifier == nil {
		i.notifier = messaging.NewDiscordNotifier(
			i.httpClient,
			props.Notifications.DiscordWebhookURL,
		)
	}

	if i.failureHandler == nil {
		i.failureHandler = messaging.NewFailureHandler(
			i.sqsHelper,
			i.notifier,
			props.Services.DLQUrl,
		)
	}

	if i.notificationPublisher == nil {
		i.notificationPublisher = publisher.NewNotificationPublisher(i.snsHelper)
	}

	// Use MaxSatt API for forest completion repository instead of database
	if i.forestCompletionRepository == nil {
		i.forestCompletionRepository = maxsattapi.NewProcessingWebService(i.maxsattAPIClient)
	}

	if i.processCompletionsService == nil {
		i.processCompletionsService = service.NewProcessCompletionsService(
			i.forestCompletionRepository,
			i.notificationPublisher,
			i.forestCompletionRepository,
		)
	}

	if i.ErrorHandler == nil {
		i.ErrorHandler = lambda.ErrorHandler(i.notifier, i.failureHandler)
	}

	if i.Handler == nil {
		i.Handler = lambda.Handler(i.ErrorHandler, i.processCompletionsService)
	}

	return i
}

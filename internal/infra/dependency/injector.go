package dependency

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/service"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"
	config "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/db"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	integrationAdapter "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/lambda"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/validator"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/persistence"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/publisher"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/utils"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/webservice"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/http"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
	redishelper "github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/redis"
)

type injector struct {
	loggerWrapper logger.Logger

	// AWS Clients
	S3             aws.S3
	Sns            aws.Sns
	DynamoDB       aws.DynamoDB
	SecretsManager aws.SecretsManager

	// External clients
	httpClient http.Client
	Redis      redis.UniversalClient

	// Helpers
	redisHelper          redishelper.RedisHelper
	s3Helper             aws.S3HelperAdapter
	snsHelper            aws.SnsHelperAdapter
	parquetProcessor     integrationAdapter.ParquetProcessor
	secretsManagerHelper aws.SecretsManagerHelperAdapter

	// Web services
	authWebService     integrationAdapter.AuthWebService
	analysisWebService integrationAdapter.AnalysisWebService
	weatherWebService  integrationAdapter.WeatherWebService
	fileWebService     integrationAdapter.FileWebService

	// Publishers
	eventPublisher integrationAdapter.EventPublisher

	// Notifiers
	notifier integrationAdapter.Notifier

	// Use cases
	processParquetInChunks  usecase.ProcessParquetInChunks
	fetchFieldAnalysis      usecase.FetchFieldAnalysis
	fetchWeatherData        usecase.FetchWeatherData
	calculateWeatherMetrics usecase.CalculateWeatherMetrics
	mergeDatasets           usecase.MergeDatasets
	fetchDeltaFile          usecase.FetchDeltaFile
	createFileRecord        usecase.CreateFileRecord
	publishEvent            usecase.PublishEvent

	// Services
	processClimateAnalysisService adapter.ProcessClimateAnalysisService

	// Entrypoint
	Validator    validator.Validate
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

func (i *injector) Wire(ctx context.Context) *injector {
	props := properties.Properties()

	if i.loggerWrapper == nil {
		i.loggerWrapper = logger.LoggerWrapper()
	}

	if i.S3 == nil {
		i.S3 = aws.S3Client(props.Aws.S3.Region)
	}

	if i.httpClient == nil {
		i.httpClient = http.HttpClient(30*time.Second, i.loggerWrapper)
	}

	if i.Redis == nil && props.Redis.Host != "" {
		i.Redis = config.Redis()
	}

	if i.redisHelper == nil {
		if i.Redis != nil {
			i.redisHelper = redishelper.Redis(i.Redis, props.Services.WeatherCacheTable, i.loggerWrapper)
		} else {
			i.redisHelper = redishelper.NoOpRedisHelper()
		}
	}

	if i.s3Helper == nil {
		i.s3Helper = aws.S3Helper(i.S3, i.loggerWrapper)
	}

	if i.Sns == nil {
		i.Sns = aws.SnsClient(props.Aws.Sns.Region)
	}

	if i.snsHelper == nil {
		i.snsHelper = aws.SnsHelper(i.Sns, i.loggerWrapper)
	}

	if i.mergeDatasets == nil {
		i.mergeDatasets = usecase.NewMergeDatasetsUseCase()
	}

	if i.parquetProcessor == nil {
		i.parquetProcessor = utils.NewParquetProcessor(i.mergeDatasets)
	}

	if i.SecretsManager == nil {
		i.SecretsManager = aws.SecretsManagerClient(props.Aws.SecretsManager.Region)
	}

	if i.secretsManagerHelper == nil {
		i.secretsManagerHelper = aws.SecretsManagerHelper(i.SecretsManager, i.redisHelper, i.loggerWrapper)
	}

	properties.Properties().Init(ctx, i.secretsManagerHelper)

	if i.authWebService == nil {
		i.authWebService = webservice.NewAuthWebService(i.httpClient, i.redisHelper)
	}

	if i.analysisWebService == nil {
		i.analysisWebService = webservice.NewAnalysisWebService(i.httpClient, i.authWebService)
	}

	if i.weatherWebService == nil {
		i.weatherWebService = webservice.NewWeather(i.httpClient, i.redisHelper)
	}

	if i.fileWebService == nil {
		i.fileWebService = webservice.NewFileWebService(i.httpClient, i.authWebService)
	}

	if i.notifier == nil {
		i.notifier = webservice.NewDiscordNotifier(
			i.httpClient,
			props.Notifications.DiscordWebhookURL,
		)
	}

	if i.eventPublisher == nil {
		i.eventPublisher = publisher.NewForestEventPublisher(i.snsHelper)
	}

	if i.processParquetInChunks == nil {
		i.processParquetInChunks = persistence.NewDeltaDataset(
			i.S3,
			i.s3Helper,
			i.parquetProcessor,
		)
	}

	if i.fetchFieldAnalysis == nil {
		i.fetchFieldAnalysis = i.analysisWebService
	}

	if i.fetchWeatherData == nil {
		i.fetchWeatherData = i.weatherWebService
	}

	if i.calculateWeatherMetrics == nil {
		i.calculateWeatherMetrics = usecase.NewCalculateWeatherMetricsUseCase()
	}

	if i.fetchDeltaFile == nil {
		i.fetchDeltaFile = i.fileWebService
	}

	if i.createFileRecord == nil {
		i.createFileRecord = i.fileWebService
	}

	if i.publishEvent == nil {
		i.publishEvent = i.eventPublisher
	}

	if i.processClimateAnalysisService == nil {
		i.processClimateAnalysisService = service.NewProcessClimateAnalysisService(
			i.processParquetInChunks,
			i.fetchFieldAnalysis,
			i.fetchWeatherData,
			i.calculateWeatherMetrics,
			i.fetchDeltaFile,
			i.createFileRecord,
			i.publishEvent,
		)
	}

	if i.Validator == nil {
		i.Validator = validator.CustomValidator()
	}

	if i.ErrorHandler == nil {
		i.ErrorHandler = lambda.ErrorHandler(i.notifier)
	}

	if i.Handler == nil {
		i.Handler = lambda.Handler(i.ErrorHandler, i.processClimateAnalysisService, i.Validator)
	}

	return i
}

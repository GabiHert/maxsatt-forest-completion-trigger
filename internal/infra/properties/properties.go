package properties

import (
	"context"
	"os"
	"strconv"
	"sync"

	aws2 "github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
)

type properties struct {
	Application   *application
	Database      *database
	Redis         *redis
	Aws           *aws
	Services      *services
	Notifications *notifications
}

type notifications struct {
	DiscordWebhookURL string
}

type application struct {
	Secrets        map[string]string
	ServerPort     string
	ServiceName    string
	Secret         string
	HealthCheckLog bool
	Profiler       bool
	Prometheus     bool
	Swagger        bool
	IsLambda       bool
}

type database struct {
	Secret   string
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	SslMode  string
	LogLevel string
	Schema   string
}

type redis struct {
	Host         string
	Port         string
	Username     string
	Password     string
	GlobalPrefix string
}

type aws struct {
	Config         *awsConfig
	SecretsManager *secretsManager
	DynamoDb       *dynamodb
	Sns            *sns
	Sqs            *sqs
	S3             *s3
	Batch          *batch
}

type awsConfig struct {
	URL string
}

type secretsManager struct {
	Region string
}

type dynamodb struct {
	Region string
}

type sns struct {
	Region string
}

type sqs struct {
	Region string
}

type s3 struct {
	Region string
}

type batch struct {
	Region string
}

type services struct {
	ForestEventsTopicArn       string
	ForestScheduleTable        string
	ForestScheduleQueueUrl     string
	ClimateDataBucket          string
	WeatherCacheTable          string
	WeatherApiUrl              string
	ForestFieldAPIURL          string
	AuthClientId               string
	AuthClientSecret           string
	AuthUrl                    string
	AuthGtwId                  string
	ForestScheduleIntervalDays int
	WeatherCacheTTLHours       int
	APITimeoutSeconds          int
	HistoricalClimateDays      int
}

var (
	initPropertiesOnce sync.Once
	secretsInitialized map[string]string
)

func (p *properties) Init(ctx context.Context, secretsManager aws2.SecretsManagerHelperAdapter) {
	initPropertiesOnce.Do(
		func() {
			if secretsManager != nil {
				err := secretsManager.GetSecret(ctx, p.Application.Secret, &secretsInitialized)
				if err != nil {
					panic("failed to get secret from secrets manager. err: " + err.Error())
				}
			}
		},
	)
}

func Properties() *properties {
	return &properties{
		Application: &application{
			ServerPort:     os.Getenv("SERVER_PORT"),
			ServiceName:    os.Getenv("SERVICE_NAME"),
			HealthCheckLog: os.Getenv("HEALTH_CHECK_LOG") == "true",
			Profiler:       os.Getenv("PROFILER") == "true",
			Prometheus:     os.Getenv("PROMETHEUS") == "true",
			Swagger:        os.Getenv("SWAGGER") == "true",
			IsLambda:       os.Getenv("LAMBDA_TASK_ROOT") != "",
			Secret:         os.Getenv("API_SECRET"),
		},
		Database: &database{
			Secret:   os.Getenv("DB_SECRET"),
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
			SslMode:  os.Getenv("DB_SSL_MODE"),
			LogLevel: os.Getenv("DB_LOG_LEVEL"),
		},
		Redis: &redis{
			Host:         os.Getenv("REDIS_HOST"),
			Port:         os.Getenv("REDIS_PORT"),
			Username:     os.Getenv("REDIS_USERNAME"),
			Password:     os.Getenv("REDIS_PASSWORD"),
			GlobalPrefix: os.Getenv("REDIS_GLOBAL_PREFIX"),
		},
		Aws: &aws{
			Config: &awsConfig{
				URL: os.Getenv("AWS_URL"),
			},
			SecretsManager: &secretsManager{
				Region: os.Getenv("AWS_SECRETS_MANAGER_REGION"),
			},
			DynamoDb: &dynamodb{
				Region: os.Getenv("AWS_DYNAMODB_REGION"),
			},
			Sns: &sns{
				Region: os.Getenv("AWS_SNS_REGION"),
			},
			Sqs: &sqs{
				Region: os.Getenv("AWS_SQS_REGION"),
			},
			S3: &s3{
				Region: os.Getenv("AWS_S3_REGION"),
			},
			Batch: &batch{
				Region: os.Getenv("AWS_BATCH_REGION"),
			},
		},
		Services: &services{
			ForestEventsTopicArn:       os.Getenv("FOREST_EVENTS_TOPIC_ARN"),
			ForestScheduleIntervalDays: parseInt(os.Getenv("FOREST_SCHEDULE_INTERVAL_DAYS"), 5),
			ForestScheduleTable:        getEnvWithDefault("FOREST_SCHEDULE_TABLE", "forests-analysis-schedules"),
			ForestScheduleQueueUrl:     os.Getenv("FOREST_SCHEDULE_QUEUE_URL"),
			ClimateDataBucket:          getEnvWithDefault("S3_BUCKET_NAME", "maxsatt-climate-data"),
			WeatherCacheTable:          getEnvWithDefault("CACHE_TABLE_NAME", "weather-cache"),
			WeatherCacheTTLHours:       parseInt(os.Getenv("WEATHER_CACHE_TTL_HOURS"), 24),
			WeatherApiUrl:              os.Getenv("WEATHER_API_URL"),
			APITimeoutSeconds:          parseInt(os.Getenv("API_TIMEOUT_SECONDS"), 30),
			ForestFieldAPIURL:          os.Getenv("FOREST_FIELD_API_URL"),
			HistoricalClimateDays:      parseInt(os.Getenv("HISTORICAL_CLIMATE_DAYS"), 30),
			AuthClientId:               secretsInitialized["client_id"],
			AuthClientSecret:           secretsInitialized["client_secret"],
			AuthUrl:                    os.Getenv("AUTH_URL"),
			AuthGtwId:                  os.Getenv("AUTH_GTW_ID"),
		},
		Notifications: &notifications{
			DiscordWebhookURL: os.Getenv("DISCORD_ERROR_WEBHOOK_URL"),
		},
	}
}

func parseInt(s string, defaultVal int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return i
}

func getEnvWithDefault(key string, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}

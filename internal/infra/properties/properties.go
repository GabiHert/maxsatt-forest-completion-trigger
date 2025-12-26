package properties

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type properties struct {
	Application   *application
	Aws           *aws
	Services      *services
	Notifications *notifications
	MaxsattAPI    *maxsattAPI
}

type notifications struct {
	DiscordWebhookURL string
}

type application struct {
	ServiceName string
	IsLambda    bool
}

type aws struct {
	Region string
	Sns    *sns
	Sqs    *sqs
}

type sns struct {
	Region string
}

type sqs struct {
	Region string
}

type services struct {
	ForestEventsTopicArn string
	DLQUrl               string
}

type maxsattAPI struct {
	BaseURL        string
	AuthURL        string
	ClientID       string
	ClientSecret   string
	TimeoutSeconds int
}

// maxsattAPISecret represents the MaxSatt API service account credentials
type maxsattAPISecret struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var (
	initPropertiesOnce sync.Once
	instance           *properties
	secretsManager     adapter.SecretsManagerAdapter
)

// InitializeSecretsManager sets the secrets manager to use for fetching credentials
// This must be called before Properties() is first invoked
func InitializeSecretsManager(sm adapter.SecretsManagerAdapter) {
	secretsManager = sm
}

// ResetProperties resets the properties singleton - should only be used in tests
// This allows tests to reinitialize properties with new environment variables
func ResetProperties() {
	instance = nil
	initPropertiesOnce = sync.Once{}
}

func Properties() *properties {
	initPropertiesOnce.Do(func() {
		instance = loadProperties()
	})
	return instance
}

func loadProperties() *properties {
	ctx := context.Background()

	// Fetch MaxSatt API credentials from Secrets Manager or environment variables
	maxsattCreds := fetchMaxsattAPICredentials(ctx, os.Getenv("MAXSATT_API_SECRET_ARN"))
	if maxsattCreds.ClientID == "" {
		maxsattCreds.ClientID = os.Getenv("MAXSATT_CLIENT_ID")
	}
	if maxsattCreds.ClientSecret == "" {
		maxsattCreds.ClientSecret = os.Getenv("MAXSATT_CLIENT_SECRET")
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = os.Getenv("AWS_SNS_REGION")
	}

	return &properties{
		Application: &application{
			ServiceName: os.Getenv("SERVICE_NAME"),
			IsLambda:    os.Getenv("LAMBDA_TASK_ROOT") != "",
		},
		Aws: &aws{
			Region: awsRegion,
			Sns: &sns{
				Region: os.Getenv("AWS_SNS_REGION"),
			},
			Sqs: &sqs{
				Region: os.Getenv("AWS_SQS_REGION"),
			},
		},
		Services: &services{
			ForestEventsTopicArn: os.Getenv("FOREST_EVENTS_TOPIC_ARN"),
			DLQUrl:               os.Getenv("DLQ_URL"),
		},
		Notifications: &notifications{
			DiscordWebhookURL: os.Getenv("DISCORD_ERROR_WEBHOOK_URL"),
		},
		MaxsattAPI: &maxsattAPI{
			BaseURL:        getEnvOrDefault("MAXSATT_API_URL", ""),
			AuthURL:        getAuthURL(),
			ClientID:       maxsattCreds.ClientID,
			ClientSecret:   maxsattCreds.ClientSecret,
			TimeoutSeconds: parseInt(os.Getenv("MAXSATT_API_TIMEOUT_SECONDS"), 30),
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

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getAuthURL returns the auth URL from MAXSATT_AUTH_URL or falls back to AUTH_URL
func getAuthURL() string {
	if val := os.Getenv("MAXSATT_AUTH_URL"); val != "" {
		return val
	}
	return os.Getenv("AUTH_URL")
}

// fetchMaxsattAPICredentials retrieves the MaxSatt API credentials from Secrets Manager
func fetchMaxsattAPICredentials(ctx context.Context, secretArn string) maxsattAPISecret {
	if secretArn == "" {
		return maxsattAPISecret{}
	}

	if secretsManager == nil {
		logger.Warn(ctx, nil, "Cannot fetch MaxSatt API credentials: secrets manager not initialized")
		return maxsattAPISecret{}
	}

	logger.Info(ctx, "Fetching MaxSatt API credentials from Secrets Manager", map[string]any{
		"secretArn": secretArn,
	})

	secretString, err := secretsManager.GetSecret(ctx, secretArn)
	if err != nil {
		logger.Error(ctx, err, "Failed to fetch MaxSatt API credentials from Secrets Manager", map[string]any{
			"secretArn": secretArn,
		})
		return maxsattAPISecret{}
	}

	var apiSecret maxsattAPISecret
	if err := json.Unmarshal([]byte(secretString), &apiSecret); err != nil {
		logger.Error(ctx, err, "Failed to parse MaxSatt API secret JSON", map[string]any{
			"secretArn": secretArn,
		})
		return maxsattAPISecret{}
	}

	logger.Info(ctx, "Successfully fetched MaxSatt API credentials from Secrets Manager")
	return apiSecret
}

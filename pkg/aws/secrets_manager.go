package aws

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerHelperAdapter interface {
	GetSecret(ctx context.Context, secret string, value any) error
}

type SecretsManager interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type secretsManagerService struct {
	secretsManager SecretsManager
	redis          redisAdapter
	logger         loggerAdapter
}

func SecretsManagerHelper(secretsManager SecretsManager, redis redisAdapter, logger loggerAdapter) SecretsManagerHelperAdapter {
	return &secretsManagerService{
		secretsManager: secretsManager,
		redis:          redis,
		logger:         logger,
	}
}

func SecretsManagerClient(region string) *secretsmanager.Client {
	cfg := getConfig(region)
	if awsUrl := os.Getenv("AWS_URL"); awsUrl != "" {
		return secretsmanager.NewFromConfig(cfg, func(o *secretsmanager.Options) {
			o.BaseEndpoint = &awsUrl
		})
	}

	return secretsmanager.NewFromConfig(cfg)
}

func (s *secretsManagerService) GetSecret(ctx context.Context, secret string, value any) error {
	s.logger.Debug(ctx, "Started", secret)

	err := s.redis.Get(ctx, secret, value)
	if err == nil {
		s.logger.Debug(ctx, "Finished", secret)
		return nil
	}

	result, err := s.secretsManager.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secret,
	})
	if err != nil {
		return err
	}

	if result == nil {
		return errors.New("secret is empty")
	}

	err = json.Unmarshal([]byte(*result.SecretString), value)
	if err != nil {
		return err
	}

	_ = s.redis.Put(ctx, secret, value)

	s.logger.Debug(ctx, "Finished", secret)
	return nil
}

package secret

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

// secretsManager implements the SecretsManagerAdapter interface
type secretsManager struct {
	helper aws.SecretsManagerHelperAdapter
	cache  map[string]string
	mu     sync.RWMutex
}

// NewSecretsManager creates a new SecretsManager that satisfies the SecretsManagerAdapter interface
func NewSecretsManager(helper aws.SecretsManagerHelperAdapter) adapter.SecretsManagerAdapter {
	return &secretsManager{
		helper: helper,
		cache:  make(map[string]string),
	}
}

// GetSecret retrieves a secret from AWS Secrets Manager with caching
func (s *secretsManager) GetSecret(ctx context.Context, secretName string) (string, error) {
	s.mu.RLock()
	if cachedValue, exists := s.cache[secretName]; exists {
		s.mu.RUnlock()
		logger.Debug(ctx, "Retrieved secret from cache", map[string]any{
			"secretName": secretName,
		})
		return cachedValue, nil
	}
	s.mu.RUnlock()

	logger.Debug(ctx, "Fetching secret from AWS Secrets Manager", map[string]any{
		"secretName": secretName,
	})

	var secretValue map[string]any
	err := s.helper.GetSecret(ctx, secretName, &secretValue)
	if err != nil {
		logger.Error(ctx, err, "Failed to retrieve secret", map[string]any{
			"secretName": secretName,
		})
		return "", err
	}

	// Marshal back to JSON string for caching
	jsonBytes, err := json.Marshal(secretValue)
	if err != nil {
		logger.Error(ctx, err, "Failed to marshal secret value", map[string]any{
			"secretName": secretName,
		})
		return "", err
	}

	secretString := string(jsonBytes)

	s.mu.Lock()
	s.cache[secretName] = secretString
	s.mu.Unlock()

	logger.Debug(ctx, "Successfully fetched and cached secret", map[string]any{
		"secretName": secretName,
	})

	return secretString, nil
}

// GetSecretJSON retrieves a secret from AWS Secrets Manager and unmarshals it into the target object
func (s *secretsManager) GetSecretJSON(ctx context.Context, secretName string, target interface{}) error {
	logger.Debug(ctx, "Fetching JSON secret from AWS Secrets Manager", map[string]any{
		"secretName": secretName,
	})

	err := s.helper.GetSecret(ctx, secretName, target)
	if err != nil {
		logger.Error(ctx, err, "Failed to retrieve or parse JSON secret", map[string]any{
			"secretName": secretName,
		})
		return err
	}

	logger.Debug(ctx, "Successfully fetched and parsed JSON secret", map[string]any{
		"secretName": secretName,
	})

	return nil
}

package adapter

import (
	"context"
)

// SecretsManagerAdapter provides methods to retrieve secrets from AWS Secrets Manager
// This adapter composes use case interfaces for secret management operations
type SecretsManagerAdapter interface {
	// GetSecret retrieves a secret as a string
	GetSecret(ctx context.Context, secretName string) (string, error)

	// GetSecretJSON retrieves a secret and unmarshals it into the target object
	GetSecretJSON(ctx context.Context, secretName string, target interface{}) error
}

package mock

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/google/uuid"
)

type SecretsManagerClient struct {
	secrets map[string]map[string]any
}

var secretsManagerMockInstance *SecretsManagerClient

func NewSecretsManagerMock() *SecretsManagerClient {
	if secretsManagerMockInstance == nil {
		secretsManagerMockInstance = &SecretsManagerClient{
			secrets: make(map[string]map[string]any),
		}
	}

	return secretsManagerMockInstance
}

func (s *SecretsManagerClient) GetSecretValue(
	_ context.Context,
	params *secretsmanager.GetSecretValueInput,
	_ ...func(*secretsmanager.Options),
) (*secretsmanager.GetSecretValueOutput, error) {
	secretName := strings.Split(*params.SecretId, ":")[len(strings.Split(*params.SecretId, ":"))-1]
	secretValue := s.secrets[secretName]

	secretBytes, err := json.Marshal(secretValue)
	if err != nil {
		return nil, err
	}

	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(string(secretBytes)),
		SecretBinary: secretBytes,
		VersionId:    aws.String(uuid.New().String()),
		ARN:          aws.String("arn:aws:secretsmanager:us-east-2:123456789012:secret:" + secretName),
		Name:         aws.String(secretName),
	}, nil
}

func (s *SecretsManagerClient) SetSecret(name string, value map[string]any) {
	s.secrets[name] = value
}

func (s *SecretsManagerClient) Reset() *SecretsManagerClient {
	s.secrets = make(map[string]map[string]any)
	return s
}

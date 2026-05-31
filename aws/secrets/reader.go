package secrets

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"
)

type SecretsManagerClient interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type SecretsReader interface {
	ReadValue(secretID string) (string, error)
}

var _ SecretsReader = (*secretsReader)(nil)

type secretsReader struct {
	client SecretsManagerClient
}

func NewSecretsReader(client SecretsManagerClient) SecretsReader {
	return &secretsReader{client: client}
}

func (r *secretsReader) ReadValue(secretID string) (string, error) {
	secret, err := r.client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretID,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to GetSecretValue")
	}
	if secret == nil || secret.SecretString == nil {
		return "", errors.New("GetSecretValue returned no SecretString")
	}
	return *secret.SecretString, nil
}

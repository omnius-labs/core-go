package secrets

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
)

type SecretsReader interface {
	ReadValue(secretID string) (string, error)
}

var _ SecretsReader = (*secretsReader)(nil)

type secretsReader struct {
	client *secretsmanager.SecretsManager
}

func NewSecretsReader(client *secretsmanager.SecretsManager) SecretsReader {
	return &secretsReader{client: client}
}

func (r *secretsReader) ReadValue(secretID string) (string, error) {
	secret, err := r.client.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &secretID,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to GetSecretValue")
	}
	return *secret.SecretString, nil
}

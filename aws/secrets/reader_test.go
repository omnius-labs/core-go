package secrets_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/omnius-labs/core-go/cloud-aws/secrets"
)

type mockSecretsManagerClient struct {
	output *secretsmanager.GetSecretValueOutput
	err    error

	gotContext  context.Context
	gotSecretID string
}

func (m *mockSecretsManagerClient) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	m.gotContext = ctx
	if params.SecretId != nil {
		m.gotSecretID = *params.SecretId
	}

	return m.output, m.err
}

func TestSecretsReaderReadValue(t *testing.T) {
	secretString := "secret-value"
	client := &mockSecretsManagerClient{
		output: &secretsmanager.GetSecretValueOutput{
			SecretString: &secretString,
		},
	}

	reader := secrets.NewSecretsReader(client)

	got, err := reader.ReadValue("secret-id")
	if err != nil {
		t.Fatalf("ReadValue returned error: %v", err)
	}
	if got != secretString {
		t.Fatalf("ReadValue = %q, want %q", got, secretString)
	}
	if client.gotContext == nil {
		t.Fatal("GetSecretValue context is nil")
	}
	if client.gotSecretID != "secret-id" {
		t.Fatalf("SecretId = %q, want %q", client.gotSecretID, "secret-id")
	}
}

func TestSecretsReaderReadValueWrapsGetSecretValueError(t *testing.T) {
	wantErr := stderrors.New("request failed")
	reader := secrets.NewSecretsReader(&mockSecretsManagerClient{
		err: wantErr,
	})

	_, err := reader.ReadValue("secret-id")
	if err == nil {
		t.Fatal("ReadValue returned nil error")
	}
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("ReadValue error = %v, want wrapping %v", err, wantErr)
	}
}

func TestSecretsReaderReadValueReturnsErrorWithoutSecretString(t *testing.T) {
	reader := secrets.NewSecretsReader(&mockSecretsManagerClient{
		output: &secretsmanager.GetSecretValueOutput{},
	})

	_, err := reader.ReadValue("secret-id")
	if err == nil {
		t.Fatal("ReadValue returned nil error")
	}
}

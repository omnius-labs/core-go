package secrets

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSecretsReader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Secrets Suite")
}

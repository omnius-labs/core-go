//go:build !ci
// +build !ci

package secrets_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSecretsReader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SecretsReader Spec")
}

var _ = Describe("SecretsReader Test", func() {
	It("simple test", func() {
		fmt.Println("This is a log message")
	})
})

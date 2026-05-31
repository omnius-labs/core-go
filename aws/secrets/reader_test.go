package secrets_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("SecretsReader Test", func() {
	It("simple test", func() {
		Skip("This is a skipped test")
		fmt.Println("This is a log message")
	})
})

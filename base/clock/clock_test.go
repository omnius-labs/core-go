package clock_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clock Spec")
}

var _ = Describe("Clock Test", func() {
	It("simple test", func() {
		fmt.Println("This is a log message")
	})
})

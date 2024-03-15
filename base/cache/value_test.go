package cache

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/omnius-labs/core-go/base/clock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestValueCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ValueCache Spec")
}

var _ = Describe("Success Test", func() {
	c := clock.NewMock(
		[]time.Time{
			time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC),
			time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC),
			time.Date(2000, time.January, 1, 1, 0, 10, 0, time.UTC),
			time.Date(2000, time.January, 1, 1, 0, 10, 0, time.UTC),
			time.Date(2000, time.January, 1, 1, 0, 60, 0, time.UTC),
		},
	)
	fr := 0
	f := func() (int, error) {
		fr++
		return fr, nil
	}

	vc := NewValueCache[int](c, 5*time.Second, 30*time.Second)
	wg := &sync.WaitGroup{}
	vc.onRefresh = func() {
		wg.Done()
	}

	It("first fetch", func() {
		wg.Add(1)
		ret, err := vc.Get(f)
		Expect(ret).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
	})

	It("use cache", func() {
		ret, err := vc.Get(f)
		Expect(ret).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
	})

	It("use cache, second fetch", func() {
		wg.Add(1)
		ret, err := vc.Get(f)
		Expect(ret).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())

		wg.Wait()

		ret, err = vc.Get(f)
		Expect(ret).To(Equal(2))
		Expect(err).NotTo(HaveOccurred())
	})

	It("third fetch", func() {
		wg.Add(1)
		ret, err := vc.Get(f)
		Expect(ret).To(Equal(3))
		Expect(err).NotTo(HaveOccurred())

		wg.Wait()
	})
})

var _ = Describe("Error Test", func() {
	c := clock.NewMock(
		[]time.Time{
			time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC),
			time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC),
			time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC),
		},
	)
	fr := 0
	f := func() (int, error) {
		fr++
		return fr, errors.New("error")
	}

	vc := NewValueCache[int](c, 5*time.Second, 30*time.Second)

	It("fetch error 1", func() {
		ret, err := vc.Get(f)
		Expect(ret).To(Equal(0))
		Expect(err).To(HaveOccurred())
		Expect(fr).To(Equal(1))
	})

	It("fetch error 2", func() {
		ret, err := vc.Get(f)
		Expect(ret).To(Equal(0))
		Expect(err).To(HaveOccurred())
		Expect(fr).To(Equal(2))
	})
})

package cache

import (
	"sync"

	"github.com/omnius-labs/core-go/base/clock"
	"golang.org/x/sync/semaphore"
)

type ValueCache[T any] struct {
	clock          clock.Clock
	data           T
	expireRefresh  int64
	expireRotten   int64
	mutex          sync.Mutex
	semaphore      *semaphore.Weighted
	timeoutRefresh int64
	timeoutRotten  int64
	onRefresh      func() // for test
}

func NewValueCache[T any](clock clock.Clock, timeoutRefresh int64, timeoutRotten int64) *ValueCache[T] {
	return &ValueCache[T]{
		clock:          clock,
		expireRefresh:  0,
		expireRotten:   0,
		mutex:          sync.Mutex{},
		semaphore:      semaphore.NewWeighted(1),
		timeoutRefresh: timeoutRefresh,
		timeoutRotten:  timeoutRotten,
	}
}

func (c *ValueCache[T]) Get(getter func() (T, error)) (T, error) {
	now := c.clock.Now().Unix()

	if now < c.expireRefresh {
		return c.data, nil
	}

	if now < c.expireRotten {
		isAcquired := c.semaphore.TryAcquire(1)
		if !isAcquired {
			return c.data, nil
		}
		data := c.data
		go func() {
			defer c.semaphore.Release(1)
			data, err := getter()
			if err != nil {
				return
			}
			c.data = data
			c.expireRefresh = now + c.timeoutRefresh
			c.expireRotten = now + c.timeoutRotten
			if c.onRefresh != nil {
				c.onRefresh()
			}
		}()
		return data, nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	data, err := getter()
	if err != nil {
		return *new(T), err
	}

	c.data = data
	c.expireRefresh = now + c.timeoutRefresh
	c.expireRotten = now + c.timeoutRotten
	if c.onRefresh != nil {
		c.onRefresh()
	}

	return data, nil
}

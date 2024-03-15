package cache

import (
	"sync"
	"time"

	"github.com/omnius-labs/core-go/base/cache/internal"
	"github.com/omnius-labs/core-go/base/clock"
	"golang.org/x/sync/semaphore"
)

type keyValuePair[T any] struct {
	key           string
	value         T
	expireRefresh int64
	expireRotten  int64
}

func newKeyValuePair[T any](key string, value T, expireRefresh int64, expireRotten int64) *keyValuePair[T] {
	return &keyValuePair[T]{
		key:           key,
		value:         value,
		expireRefresh: expireRefresh,
		expireRotten:  expireRotten,
	}
}

type KeyValueCache[T any] struct {
	clock          clock.Clock
	dict           *internal.SyncMap[string, *internal.LinkedListNode[*keyValuePair[T]]]
	keys           *internal.LinkedList[*keyValuePair[T]]
	capacity       int
	mutex          sync.Mutex
	semaphore      *semaphore.Weighted
	timeoutRefresh int64
	timeoutRotten  int64
	onRefresh      func() // for test
}

func NewKeyValueCache[T any](clock clock.Clock, capacity int, timeoutRefresh time.Duration, timeoutRotten time.Duration) *KeyValueCache[T] {
	return &KeyValueCache[T]{
		clock:          clock,
		dict:           internal.NewSyncMap[string, *internal.LinkedListNode[*keyValuePair[T]]](),
		keys:           internal.NewLinkedList[*keyValuePair[T]](),
		capacity:       capacity,
		mutex:          sync.Mutex{},
		semaphore:      semaphore.NewWeighted(1),
		timeoutRefresh: int64(timeoutRefresh.Seconds()),
		timeoutRotten:  int64(timeoutRotten.Seconds()),
	}
}

func (c *KeyValueCache[T]) Get(key string, getter func() (T, error)) (T, error) {
	now := c.clock.Now().Unix()

	node, ok := c.dict.Get(key)

	if ok && now < node.Value.expireRefresh {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		c.keys.Remove(node)
		c.keys.AppendLast(node)

		return node.Value.value, nil
	}

	if ok && now < node.Value.expireRotten {
		isAcquired := c.semaphore.TryAcquire(1)
		if !isAcquired {
			return node.Value.value, nil
		}
		value := node.Value.value
		go func() {
			defer c.semaphore.Release(1)
			value, err := getter()
			if err != nil {
				return
			}
			node.Value.value = value
			node.Value.expireRefresh = now + c.timeoutRefresh
			node.Value.expireRotten = now + c.timeoutRotten
			if c.onRefresh != nil {
				c.onRefresh()
			}
		}()
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.keys.Remove(node)
		c.keys.AppendLast(node)
		return value, nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	data, err := getter()
	if err != nil {
		return *new(T), err
	}

	if ok {
		c.keys.Remove(node)
	}

	if c.keys.Len() >= c.capacity {
		node := c.keys.First()
		c.dict.Delete(node.Value.key)
		c.keys.Remove(node)
	}

	pair := newKeyValuePair[T](key, data, now+c.timeoutRefresh, now+c.timeoutRotten)
	node = internal.NewLinkedListNode[*keyValuePair[T]](pair)
	c.keys.AppendLast(node)
	c.dict.Set(key, node)
	if c.onRefresh != nil {
		c.onRefresh()
	}
	return data, nil
}

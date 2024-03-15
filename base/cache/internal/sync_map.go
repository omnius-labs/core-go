package internal

import "sync"

type SyncMap[TKey comparable, T any] struct {
	dict  map[TKey]T
	mutex sync.Mutex
}

func NewSyncMap[TKey comparable, T any]() *SyncMap[TKey, T] {
	return &SyncMap[TKey, T]{dict: make(map[TKey]T)}
}

func (m *SyncMap[TKey, T]) Get(key TKey) (T, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	value, ok := m.dict[key]
	return value, ok
}

func (m *SyncMap[TKey, T]) Set(key TKey, value T) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.dict[key] = value
}

func (m *SyncMap[TKey, T]) Delete(key TKey) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.dict, key)
}

func (m *SyncMap[TKey, T]) Len() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return len(m.dict)
}

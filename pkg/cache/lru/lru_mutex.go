package lru

import (
	"sync"
)

type LRUMutex[K comparable, V any] struct {
	cap  int
	lru  *LRU[K, V]
	lock sync.RWMutex
}

func NewLRUMutex[K comparable, V any](cap int, onEvicted func(key K, value V)) *LRUMutex[K, V] {
	return &LRUMutex[K, V]{
		cap:  cap,
		lru:  NewLRU[K, V](cap, onEvicted),
		lock: sync.RWMutex{},
	}
}

func (m *LRUMutex[K, V]) Remove(key K) (k K, v V, ok bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.lru == nil {
		return
	}

	return m.lru.Remove(key)
}

func (m *LRUMutex[K, V]) Contains(key K) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.lru == nil {
		return false
	}

	return m.lru.Contains(key)
}

func (m *LRUMutex[K, V]) Keys() []K {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.lru == nil {
		return nil
	}

	return m.lru.Keys()
}

func (m *LRUMutex[K, V]) Values() []V {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.lru == nil {
		return nil
	}

	return m.lru.Values()
}

func (m *LRUMutex[K, V]) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.lru == nil {
		return 0
	}

	return m.lru.Len()
}

func (m *LRUMutex[K, V]) Cap() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.lru == nil {
		return 0
	}

	return m.lru.Cap()
}

func (m *LRUMutex[K, V]) Add(key K, value V) bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.lru == nil {
		m.lru = NewLRU[K, V](m.Cap(), nil)
	}

	return m.lru.Add(key, value)
}

func (m *LRUMutex[K, V]) Get(key K) (value V, ok bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if m.lru == nil {
		return
	}

	if v, ok := m.lru.Get(key); ok {
		return v, ok
	}

	return
}

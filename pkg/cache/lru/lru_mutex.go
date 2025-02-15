package lru

import "sync"

type Mutex[K comparable, V any] struct {
	Cap  int
	lru  *LRU[K, V]
	lock sync.Mutex
}

func (m *Mutex[K, V]) Add(key K, value V) bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.lru == nil {
		m.lru = NewLRU[K, V](m.Cap, nil)
	}

	return m.lru.Add(key, value)
}

func (m *Mutex[K, V]) Get(key K) (value V, ok bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.lru == nil {
		return
	}

	if v, ok := m.lru.Get(key); ok {
		return v, ok
	}

	return
}

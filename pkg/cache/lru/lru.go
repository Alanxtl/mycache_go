package lru

import (
	"container/list"
)

type LRU[K comparable, V any] struct {
	cap       int
	cache     map[K]*list.Element
	list      *list.List
	OnEvicted func(key K, value V)
}

type Element[K comparable, V any] struct {
	key   K
	value V
}

func NewElement[K comparable, V any](key K, value V) *Element[K, V] {
	return &Element[K, V]{
		key:   key,
		value: value,
	}
}

func NewLRU[K comparable, V any](cap int, onEvicted func(key K, value V)) *LRU[K, V] {
	return &LRU[K, V]{
		cap:       cap,
		cache:     make(map[K]*list.Element, cap),
		list:      list.New(),
		OnEvicted: onEvicted,
	}
}

func (c LRU[K, V]) Add(key K, value V) bool {
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*Element[K, V]).value = value
		return false
	}

	elem := c.list.PushFront(NewElement(key, value))
	c.cache[key] = elem

	evict := len(c.cache) > c.Cap()

	if evict {
		c.RemoveOldest()
	}

	return evict
}

func (c LRU[K, V]) Remove(key K) (k K, v V, ok bool) {
	if elem, ok := c.cache[key]; ok {
		c.list.Remove(elem)
		delete(c.cache, elem.Value.(*Element[K, V]).key)

		if c.OnEvicted != nil {
			c.OnEvicted(elem.Value.(*Element[K, V]).key, elem.Value.(*Element[K, V]).value)
		}

		return elem.Value.(*Element[K, V]).key, elem.Value.(*Element[K, V]).value, true
	}
	return
}

func (c LRU[K, V]) RemoveOldest() (k K, v V, ok bool) {
	if elem := c.list.Back(); elem != nil {
		c.list.Remove(elem)
		delete(c.cache, elem.Value.(*Element[K, V]).key)

		if c.OnEvicted != nil {
			c.OnEvicted(elem.Value.(*Element[K, V]).key, elem.Value.(*Element[K, V]).value)
		}

		return elem.Value.(*Element[K, V]).key, elem.Value.(*Element[K, V]).value, true
	}
	return
}

func (c LRU[K, V]) Get(key K) (value V, ok bool) {
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*Element[K, V]).value, true
	}
	return
}

func (c LRU[K, V]) Contains(key K) bool {
	if _, ok := c.cache[key]; ok {
		return true
	}
	return false
}

func (c LRU[K, V]) Keys() []K {
	values := make([]K, 0, c.Len())

	for v := c.list.Front(); v != nil; v = v.Next() {
		values = append(values, v.Value.(*Element[K, V]).key)
	}

	return values
}

func (c LRU[K, V]) Values() []V {
	values := make([]V, 0, c.Len())

	for v := c.list.Front(); v != nil; v = v.Next() {
		values = append(values, v.Value.(*Element[K, V]).value)
	}

	return values
}

func (c LRU[K, V]) Len() int {
	return len(c.cache)
}

func (c LRU[K, V]) Cap() int {
	return c.cap
}

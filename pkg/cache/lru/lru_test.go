// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lru

import (
	"reflect"
	"testing"
)

func TestLRU(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) {
		if k != v {
			t.Fatalf("Evict values not equal (%v!=%v)", k, v)
		}
		evictCounter++
	}
	l := NewLRU(128, onEvicted)

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	if l.Len() != 128 {
		t.Fatalf("bad len: %v", l.Len())
	}
	if l.Cap() != 128 {
		t.Fatalf("expect %d, but %d", 128, l.Cap())
	}

	if evictCounter != 128 {
		t.Fatalf("bad evict count: %v", evictCounter)
	}

	for i, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || v != k || v != 256-1-i {
			t.Fatalf("bad key: %v", v)
		}
	}
	for i, v := range l.Values() {
		if v != i+128 {
			t.Fatalf("bad value: %v", v)
		}
	}
	for i := 0; i < 128; i++ {
		if _, ok := l.Get(i); ok {
			t.Fatalf("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		if _, ok := l.Get(i); !ok {
			t.Fatalf("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		if _, _, ok := l.Remove(i); !ok {
			t.Fatalf("should be contained")
		}
		if _, _, ok := l.Remove(i); ok {
			t.Fatalf("should not be contained")
		}
		if _, ok := l.Get(i); ok {
			t.Fatalf("should be deleted")
		}
	}

	l.Get(192) // expect 192 to be last key in l.Keys()

	for i, k := range l.Keys() {
		if (256-1-i < 63 && k != 256-1-i+193) || (256-1-i == 63 && k != 192) {
			t.Fatalf("out of order key: %v %v", i, k)
		}
	}

}

func TestLRU_GetOldest_RemoveOldest(t *testing.T) {
	l := NewLRU[int, int](128, nil)

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}

	k, _, ok := l.RemoveOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 128 {
		t.Fatalf("bad: %v", k)
	}

	k, _, ok = l.RemoveOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 129 {
		t.Fatalf("bad: %v", k)
	}
}

// Test that Add returns true/false if an eviction occurred
func TestLRU_Add(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) {
		evictCounter++
	}

	l := NewLRU(1, onEvicted)

	if l.Add(1, 1) == true || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	if l.Add(2, 2) == false || evictCounter != 1 {
		t.Errorf("should have an eviction")
	}
}

// Test that Contains doesn't update recent-ness
func TestLRU_Contains(t *testing.T) {
	l := NewLRU[int, int](2, nil)

	l.Add(1, 1)
	l.Add(2, 2)
	if !l.Contains(1) {
		t.Errorf("1 should be contained")
	}

	l.Add(3, 3)
	if l.Contains(1) {
		t.Errorf("Contains should not have updated recent-ness of 1")
	}
}

// Test that Peek doesn't update recent-ness
func TestLRU_Peek(t *testing.T) {
	l := NewLRU[int, int](2, nil)

	l.Add(1, 1)
	l.Add(2, 2)

	l.Add(3, 3)
	if l.Contains(1) {
		t.Errorf("should not have updated recent-ness of 1")
	}
}

func (c *LRU[K, V]) wantKeys(t *testing.T, want []K) {
	t.Helper()
	got := c.Keys()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wrong keys got: %v, want: %v ", got, want)
	}
}

func TestCache_EvictionSameKey(t *testing.T) {
	var evictedKeys []int

	cache := NewLRU(
		2,
		func(key int, _ struct{}) {
			evictedKeys = append(evictedKeys, key)
		})

	if evicted := cache.Add(1, struct{}{}); evicted {
		t.Error("First 1: got unexpected eviction")
	}
	cache.wantKeys(t, []int{1})

	if evicted := cache.Add(2, struct{}{}); evicted {
		t.Error("2: got unexpected eviction")
	}
	cache.wantKeys(t, []int{2, 1})

	if evicted := cache.Add(1, struct{}{}); evicted {
		t.Error("Second 1: got unexpected eviction")
	}
	cache.wantKeys(t, []int{1, 2})

	if evicted := cache.Add(3, struct{}{}); !evicted {
		t.Error("3: did not get expected eviction")
	}
	cache.wantKeys(t, []int{3, 1})

	want := []int{2}
	if !reflect.DeepEqual(evictedKeys, want) {
		t.Errorf("evictedKeys got: %v want: %v", evictedKeys, want)
	}
}

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := NewLRU[String, String](2, nil)
	lru.Add(String("key1"), String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	lru := NewLRU[string, String](3, nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, _, ok := lru.Remove("key1"); !ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed %v %v", ok, lru.Len())
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value String) {
		keys = append(keys, key)
	}
	lru := NewLRU[string, String](2, callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s %s", expect, keys)
	}
}

func TestAdd(t *testing.T) {
	lru := NewLRU[string, String](6, nil)
	lru.Add("key", String("1"))
	lru.Add("key", String("111"))

	if lru.Len() != 1 {
		t.Fatal("expected 6 but got", lru.Len())
	}
}

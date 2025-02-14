// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cache

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
	cache.wantKeys(t, []int{1, 2})

	if evicted := cache.Add(1, struct{}{}); evicted {
		t.Error("Second 1: got unexpected eviction")
	}
	cache.wantKeys(t, []int{2, 1})

	if evicted := cache.Add(3, struct{}{}); !evicted {
		t.Error("3: did not get expected eviction")
	}
	cache.wantKeys(t, []int{1, 3})

	want := []int{2}
	if !reflect.DeepEqual(evictedKeys, want) {
		t.Errorf("evictedKeys got: %v want: %v", evictedKeys, want)
	}
}

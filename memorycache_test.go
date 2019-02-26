package gocache

import (
	"testing"
	"time"
)

type Animal struct {
	name string
}

func TestNewMemoryCache(t *testing.T) {
	memory := NewMemoryCache()

	if memory == nil {
		t.Error("map is null")
	}
}

func TestMemoryCache_Add(t *testing.T) {
	memory := NewMemoryCache()
	elephant := Animal{"elephant"}
	monkey := Animal{"monkey"}

	memory.Add("elephant", elephant)

	policy := NewCacheItemPolicy()
	policy.SlidingExpiration = time.Second * 20

	memory.AddWithPolicy("monkey", monkey, policy)

	if memory.Count() != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestMemoryCache_Get(t *testing.T) {
	memory := NewMemoryCache()

	if _, ok := memory.Get("elephant"); ok {
		t.Error("ok should be false when item is missing from map.")
	}
	elephant := Animal{"elephant"}
	memory.Add("elephant", elephant)

	if v, ok := memory.Get("elephant"); !ok {
		t.Error("ok should be true for item stored within the map.")
	} else if elephant, ok := v.(Animal); !ok {
		t.Error("ok should be true for item stored within the map.")
	} else if elephant.name != "elephant" {
		t.Error("item was modified.")
	}
}

func TestMemoryCache_Contains(t *testing.T) {
	memory := NewMemoryCache()

	if ok := memory.ContainsKey("elephant"); ok {
		t.Error("ok should be false when item is missing from map.")
	}
	elephant := Animal{"elephant"}
	memory.Add("elephant", elephant)

	if ok := memory.ContainsKey("elephant"); !ok {
		t.Error("ok should be true when item is missing from map.")
	}
}

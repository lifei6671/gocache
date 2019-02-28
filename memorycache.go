package gocache

import (
	"context"
	"fmt"
	"time"
)

var SHARD_COUNT = 32

type MemoryCache struct {
	stores     []*MemoryCacheStore
	storeCount int
	closed     bool
	cancel     func()
	interval   time.Duration
}

func NewMemoryCache(interval time.Duration) *MemoryCache {

	cache := &MemoryCache{
		stores:     make([]*MemoryCacheStore, SHARD_COUNT),
		storeCount: SHARD_COUNT,
		interval:   interval,
	}
	ctx, cancel := context.WithCancel(context.Background())

	cache.cancel = cancel

	for i := 0; i < SHARD_COUNT; i++ {
		cache.stores[i] = NewMemoryCacheStore(ctx, fmt.Sprintf("REGION-%d", i), interval)
	}

	return cache
}

func (memory *MemoryCache) Add(key string, value interface{}) {
	store := memory.getStore(key)
	item := NewCacheItem(key, value)
	store.Add(item)
}

func (memory *MemoryCache) AddWithCacheItem(item CacheItem) {
	store := memory.getStore(item.Key)
	store.Add(&item)
}

// 增加缓存并设置策略
func (memory *MemoryCache) AddWithSlidingExpiration(key string, value interface{}, expiration time.Duration) {
	store := memory.getStore(key)
	store.AddWithSlidingExpiration(key, value, expiration)
}

func (memory *MemoryCache) AddWithAbsoluteExpiration(key string, value interface{}, expiration time.Time) {
	store := memory.getStore(key)
	store.AddWithAbsoluteExpiration(key, value, expiration)
}

func (memory *MemoryCache) Get(key string) (value interface{}, ok bool) {
	store := memory.getStore(key)
	return store.GetValue(key)
}

func (memory *MemoryCache) Remove(key string) (value interface{}, ok bool) {
	store := memory.getStore(key)
	return store.Remove(key)
}

// ContainsKey 判断缓存中是否存在指定键
func (memory *MemoryCache) ContainsKey(key string) bool {
	store := memory.getStore(key)
	return store.ContainsKey(key)
}

func (memory *MemoryCache) Count() int {
	c := 0
	for _, store := range memory.stores {
		c += store.Count()
	}
	return c
}

func (memory *MemoryCache) Clear() {
	for _, store := range memory.stores {
		store.Clear()
	}
}

func (memory *MemoryCache) Close() {
	memory.closed = true
	for _, store := range memory.stores {
		store.Close()
	}
}

func (memory *MemoryCache) getStore(key string) *MemoryCacheStore {
	if memory.closed {
		panic("cache closed")
	}
	idx := int(hashKey(key) % uint32(memory.storeCount))

	return memory.stores[idx]
}

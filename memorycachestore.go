package gocache

import (
	"context"
	"sync"
	"time"
)

type MemoryCacheStore struct {
	sync.RWMutex
	// 缓存容器
	entries map[string]*CacheItem
	// 存在滑动过期时间的缓存
	slidingExpires map[string]*CacheItem
	//存在绝对过期时间
	absoluteExpires map[string]*CacheItem
	//扫描间隔
	scanInterval time.Duration
	once         sync.Once
	cancel       func()
}

func NewMemoryCacheStore(ctx context.Context, interval time.Duration) *MemoryCacheStore {
	store := &MemoryCacheStore{
		entries:         make(map[string]*CacheItem),
		slidingExpires:  make(map[string]*CacheItem),
		absoluteExpires: make(map[string]*CacheItem),
		scanInterval:    interval,
	}
	ctx, cancel := context.WithCancel(ctx)

	store.cancel = cancel

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()

	return store
}

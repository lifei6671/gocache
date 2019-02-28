package gocache

import (
	"context"
	"log"
	"sync"
	"time"
)

type MemoryCacheStore struct {
	sync.RWMutex
	regionName string
	// 缓存容器
	entries map[string]*CacheItem
	// 存在过期时间的缓存
	expires map[string]*CacheItem
	//扫描间隔
	scanInterval time.Duration
	once         sync.Once
	cancel       func()
}

func NewMemoryCacheStore(ctx context.Context, regionName string, interval time.Duration) *MemoryCacheStore {
	if interval <= 0 {
		interval = time.Second * 1
	}
	store := &MemoryCacheStore{
		regionName:   regionName,
		entries:      make(map[string]*CacheItem, 2000),
		expires:      make(map[string]*CacheItem, 2000),
		scanInterval: interval,
	}
	ctx, cancel := context.WithCancel(ctx)

	store.cancel = cancel

	go func() {
		timer := time.NewTimer(store.scanInterval)
		for {
			select {
			case <-timer.C:
				store.RLock()
				for _, item := range store.expires {
					store.RUnlock()
					store.check(item)
					store.RLock()
				}
				store.RUnlock()

				timer.Reset(store.scanInterval)
			case <-ctx.Done():
				log.Println("缓存过期扫码已停止 ->", store.regionName)
				return
			}
		}
	}()

	return store
}

func (store *MemoryCacheStore) Add(item *CacheItem) {
	store.RLock()
	existingItem, ok := store.entries[item.Key]
	store.RUnlock()

	if ok {
		existingItem.callRemovedCallback(CacheEntryRemovedReasonCacheSpecificEviction)
		store.Lock()
		store.entries[item.Key] = item

		if item.HasExpiration() {
			store.expires[item.Key] = item
		} else {
			delete(store.expires, item.Key)
		}
		store.Unlock()

	} else {
		store.Lock()
		store.entries[item.Key] = item
		//如果存在过期时间
		if item.HasExpiration() {
			store.expires[item.Key] = item
		}
		store.Unlock()
	}
}

func (store *MemoryCacheStore) AddWithSlidingExpiration(key string, value interface{}, expiration time.Duration) {
	item := NewCacheItemWithSlidingExpiration(key, value, expiration)

	store.Add(item)
}

func (store *MemoryCacheStore) AddWithAbsoluteExpiration(key string, value interface{}, expiration time.Time) {
	item := NewCacheItemWithAbsoluteExpiration(key, value, expiration)

	store.Add(item)
}

func (store *MemoryCacheStore) GetValue(key string) (value interface{}, ok bool) {
	store.RLock()
	existingItem, ok := store.entries[key]
	//log.Println(existingItem.Key)
	store.RUnlock()
	if ok {
		existingItem.Lock()
		defer existingItem.Unlock()
		if !existingItem.InExpires() {
			existingItem.keepLive()
			return existingItem.Value, ok
		}
		//如果已过期，且存在滑动过期时间，且存在创建缓存方法
		if existingItem.InExpires() && existingItem.slidingExpiration > 0 && existingItem.CreateCallback != nil {
			if v, err := existingItem.CreateCallback(key, existingItem.Value); err != nil {
				store.remove(store.expires[key], CacheEntryRemovedReasonCacheSpecificEviction)
				return nil, false
			} else {
				existingItem.Value = v
				existingItem.keepLive()
				return v, true
			}
		}

		store.remove(existingItem, CacheEntryRemovedReasonCacheSpecificEviction)
		return nil, false

	}
	return nil, false
}

func (store *MemoryCacheStore) ContainsKey(key string) bool {
	_, ok := store.GetValue(key)
	return ok
}

func (store *MemoryCacheStore) Count() int {
	store.RLock()
	c := len(store.entries)
	store.RUnlock()
	return c
}

func (store *MemoryCacheStore) check(item *CacheItem) {
	item.Lock()
	defer item.Unlock()
	if !item.InExpires() {
		return
	}
	//如果没有设置滑动过期时间，或者没有设置更新缓存的方法，则直接删除缓存项
	if item.slidingExpiration <= 0 || item.CreateCallback == nil {
		store.Lock()
		delete(store.entries, item.Key)
		delete(store.expires, item.Key)

		store.Unlock()
		item.callRemovedCallback(CacheEntryRemovedReasonExpired)
		return
	}

	if v, err := item.CreateCallback(item.Key, item.Value); err != nil {
		log.Printf("[%s] 创建缓存值失败 -> %s - %s", store.regionName, item.Key, err)
		item.callRemovedCallback(CacheEntryRemovedReasonCacheSpecificEviction)
		store.Lock()
		delete(store.entries, item.Key)
		delete(store.expires, item.Key)
		store.Unlock()
	} else {
		item.Value = v
		item.keepLive()
	}
}

func (store *MemoryCacheStore) remove(item *CacheItem, reason CacheEntryRemovedReason) {
	if item == nil {
		return
	}
	store.Lock()
	delete(store.expires, item.Key)
	delete(store.entries, item.Key)
	if item.RemovedCallback != nil {
		item.RemovedCallback(item.Key, item.Value, reason)
	}
	store.Unlock()
}

package gocache

import (
	"log"
	"sync"
	"time"
)

var MaxTimeValue = time.Date(9999, 0, 0, 0, 0, 0, 0, time.Local)

type MemoryCacheStore struct {
	// 是否释放资源
	closed bool
	sync.RWMutex
	// 缓存容器
	entries map[string]*memoryCacheEntry
	// 存在过期时间的缓存
	expires map[string]*memoryCacheEntry
}

func (store *MemoryCacheStore) Add(key string, value interface{}) {

	store.AddWithPolicy(key, value, NewCacheItemPolicy())
}

// AddWithPolicy 添加缓存并设置过期策略
// 如果设置的缓存已存在，则更新缓存的值，并重新设置过期策略，同时发送更新回调
// 如果设置的缓存不存在，则创建一个缓存值，更加入到缓存中
func (store *MemoryCacheStore) AddWithPolicy(key string, value interface{}, policy *CacheItemPolicy) {
	var existingEntry *memoryCacheEntry

	store.RLock()
	// 如果已存在
	if existingEntry, _ = store.entries[key]; existingEntry != nil {
		existingEntry.value = value
		*(existingEntry.policy) = *policy
		store.RUnlock()
		store.update(existingEntry, CacheEntryRemovedReasonRemoved)
	} else {
		store.RUnlock()
		if policy.SlidingExpiration > 0 {
			policy.AbsoluteExpiration = time.Now().Add(policy.SlidingExpiration)
		}
		existingEntry = newMemoryCacheEntry()
		existingEntry.value = value
		existingEntry.Key = key
		existingEntry.policy = policy

		store.Lock()
		store.entries[key] = existingEntry
		store.Unlock()
	}
}

func (store *MemoryCacheStore) Get(key string) (value interface{}, ok bool) {
	var existingEntry *memoryCacheEntry
	added := false
	store.RLock()
	existingEntry, added = store.entries[key]
	store.RUnlock()
	ok = added
	if ok {
		value = existingEntry.value
	}
	return
}

func (store *MemoryCacheStore) Count() int {
	store.RLock()
	c := len(store.entries)
	store.RUnlock()
	return c
}

// 判断一个键是否存在
func (store *MemoryCacheStore) ContainsKey(key string) bool {

	var existingEntry *memoryCacheEntry
	ok := false

	store.RLock()
	//如果已存在并且没有过期
	if existingEntry, ok = store.entries[key]; !ok {
		store.RUnlock()
		return false
	}
	return true

	//如果存在且未过期
	if !existingEntry.isExpired() {
		store.RUnlock()
		return true
	}
	log.Println("aaaa")

	if existingEntry.policy.CreateCallback != nil {

		if vv, err := existingEntry.policy.CreateCallback(existingEntry.Key); err != nil {
			store.RUnlock()
			store.remove(existingEntry, CacheEntryRemovedReasonExpired)
			return false
		} else {
			existingEntry.value = vv
			store.RUnlock()
			store.update(existingEntry, CacheEntryRemovedReasonExpired)
			return true
		}
	}
	store.RUnlock()
	store.remove(existingEntry, CacheEntryRemovedReasonExpired)
	return false
}

func (store *MemoryCacheStore) remove(entry *memoryCacheEntry, reason CacheEntryRemovedReason) {
	if entry == nil {
		return
	}
	store.Lock()
	delete(store.entries, entry.Key)
	store.Unlock()
	if entry.policy.RemovedCallback != nil {
		args := CacheEntryRemovedArguments{
			RemovedReason: reason,
			CacheItem:     CacheItem{Key: entry.Key, Value: entry.value},
		}
		entry.policy.RemovedCallback(args)
	}
}

func (store *MemoryCacheStore) update(entry *memoryCacheEntry, reason CacheEntryRemovedReason) {
	if entry == nil {
		return
	}
	store.Lock()
	if entry.policy.UpdateCallback != nil {
		args := CacheEntryUpdateArguments{
			RemovedReason:    reason,
			UpdatedCacheItem: CacheItem{Key: entry.Key, Value: entry.value},
		}
		entry.policy.UpdateCallback(args)
	}
	if entry.policy.SlidingExpiration > 0 {
		entry.policy.AbsoluteExpiration = time.Now().Add(entry.policy.SlidingExpiration)
	}
	store.Unlock()
}

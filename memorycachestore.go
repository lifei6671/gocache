package gocache

import (
	"context"
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
	once    sync.Once
	cancel  func()
	ctx     context.Context
}

func NewMemoryCacheStore(ctx context.Context) *MemoryCacheStore {
	store := &MemoryCacheStore{
		closed:  false,
		entries: make(map[string]*memoryCacheEntry, 1000),
		expires: make(map[string]*memoryCacheEntry, 1000),
	}
	ctx, cancel := context.WithCancel(ctx)

	store.cancel = cancel
	store.ctx = ctx

	return store
}

func (store *MemoryCacheStore) Add(key string, value interface{}) {

	policy := NewCacheItemPolicy()
	store.AddWithPolicy(key, value, &policy)
}

// AddWithPolicy 添加缓存并设置过期策略
// 如果设置的缓存已存在，则更新缓存的值，并重新设置过期策略，同时发送更新回调
// 如果设置的缓存不存在，则创建一个缓存值，更加入到缓存中
func (store *MemoryCacheStore) AddWithPolicy(key string, value interface{}, policy *CacheItemPolicy) {

	store.once.Do(func() {
		go func() {
			timer := time.NewTimer(time.Second * 1)
			defer timer.Stop()

			for {
				select {
				case <-timer.C:
					timer.Stop()
					store.Lock()
					for _, entry := range store.expires {
						entry.Lock()
						if entry.isExpired() {
							store.removeFromCache(entry, CacheEntryRemovedReasonExpired)
						}
						entry.Unlock()
					}
					store.Unlock()
					timer.Reset(time.Second * 1)
				case <-store.ctx.Done():
					return
				}
			}
		}()
	})

	var existingEntry *memoryCacheEntry

	store.RLock()
	// 如果已存在
	if existingEntry, _ = store.entries[key]; existingEntry != nil {
		existingEntry.value = value
		*(existingEntry.policy) = *policy
		store.RUnlock()
		existingEntry.Lock()
		store.update(existingEntry, CacheEntryRemovedReasonRemoved)
		existingEntry.Unlock()
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
		if existingEntry.hasExpiration() {
			store.expires[key] = existingEntry
		}
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
	if existingEntry != nil {
		existingEntry.Lock()
		defer existingEntry.Unlock()
	}
	if ok && !existingEntry.isExpired() {
		value = existingEntry.value
		existingEntry.keep()
	}
	if existingEntry != nil && existingEntry.isExpired() {
		if existingEntry.policy.CreateCallback != nil {
			if v, err := existingEntry.policy.CreateCallback(key, existingEntry.value); err == nil {
				existingEntry.value = v
				store.update(existingEntry, CacheEntryRemovedReasonExpired)
				return v, true
			}
		}
		store.Lock()
		store.remove(existingEntry, CacheEntryRemovedReasonExpired)
		store.Unlock()
		ok = false
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

		return false
	}
	existingEntry.Lock()
	defer existingEntry.Unlock()
	store.RUnlock()

	//如果存在且未过期
	if !existingEntry.isExpired() {
		existingEntry.keep()
		return true
	}

	if existingEntry.policy.CreateCallback != nil {

		if vv, err := existingEntry.policy.CreateCallback(existingEntry.Key, existingEntry.value); err != nil {
			store.Lock()
			store.remove(existingEntry, CacheEntryRemovedReasonExpired)
			store.Unlock()
			return false
		} else {
			existingEntry.value = vv
			store.update(existingEntry, CacheEntryRemovedReasonExpired)
			return true
		}
	}

	store.Lock()
	store.remove(existingEntry, CacheEntryRemovedReasonExpired)
	store.Unlock()

	return false
}

func (store *MemoryCacheStore) check(entry *memoryCacheEntry) {
	if entry == nil {
		return
	}
	//1、对缓存项加锁，防止其他goroutine操作
	//2、检查是否过期，且是否存在创建缓存的回调函数
	//3、如果创建缓存函数，则调用生成新的缓存，如果出错则删除缓存，并判断实付存在删除通知回调
	//4、如果不存在创建函数，则删除缓存，并调用删除回调方法
	//5、函数未对上级容器加锁，因此需要在调用方加锁
	entry.Lock()
	if entry.isExpired() && entry.policy.CreateCallback != nil {
		if v, err := entry.policy.CreateCallback(entry.Key, entry.value); err != nil {
			if entry.policy.RemovedCallback != nil {
				args := CacheEntryRemovedArguments{
					CacheItem:     CacheItem{Key: entry.Key, Value: entry.value},
					RemovedReason: CacheEntryRemovedReasonCacheSpecificEviction,
				}
				entry.policy.RemovedCallback(args)
			}
			delete(store.expires, entry.Key)
			delete(store.entries, entry.Key)
		} else {
			entry.value = v
			if entry.policy.SlidingExpiration > 0 {
				entry.policy.AbsoluteExpiration = time.Now().Add(entry.policy.SlidingExpiration)
			}
		}
	} else if entry.isExpired() {
		delete(store.expires, entry.Key)
		delete(store.entries, entry.Key)
		args := CacheEntryRemovedArguments{
			CacheItem:     CacheItem{Key: entry.Key, Value: entry.value},
			RemovedReason: CacheEntryRemovedReasonExpired,
		}
		entry.policy.RemovedCallback(args)
	}
	entry.Unlock()

}

func (store *MemoryCacheStore) removeFromCache(entry *memoryCacheEntry, reason CacheEntryRemovedReason) {
	if entry == nil {
		return
	}
	delete(store.entries, entry.Key)
	delete(store.expires, entry.Key)
	if entry.policy.RemovedCallback != nil {
		args := CacheEntryRemovedArguments{
			RemovedReason: reason,
			CacheItem:     CacheItem{Key: entry.Key, Value: entry.value},
		}
		entry.policy.RemovedCallback(args)
	}
}

func (store *MemoryCacheStore) remove(entry *memoryCacheEntry, reason CacheEntryRemovedReason) {
	if entry == nil {
		return
	}
	if entry.isExpired() && entry.policy.CreateCallback == nil {
		delete(store.entries, entry.Key)
		delete(store.expires, entry.Key)
		if entry.policy.RemovedCallback != nil {
			args := CacheEntryRemovedArguments{
				RemovedReason: reason,
				CacheItem:     CacheItem{Key: entry.Key, Value: entry.value},
			}
			entry.policy.RemovedCallback(args)
		}
	}
}

func (store *MemoryCacheStore) update(entry *memoryCacheEntry, reason CacheEntryRemovedReason) {
	if entry == nil {
		return
	}
	if entry.policy.UpdateCallback != nil {
		args := CacheEntryUpdateArguments{
			RemovedReason:    reason,
			UpdatedCacheItem: CacheItem{Key: entry.Key, Value: entry.value},
		}
		entry.policy.UpdateCallback(args)
	}
	entry.keep()
}

func (store *MemoryCacheStore) Close() {
	store.Lock()
	defer store.Unlock()
	store.closed = true
	if store.cancel != nil {
		store.cancel()
	}
	store.entries = nil
	store.expires = nil

}

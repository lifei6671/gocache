package gocache

import (
	"sync"
	"time"
)

type memoryCacheEntry struct {
	Key     string
	value   interface{}
	created time.Time
	policy  *CacheItemPolicy
	sync.RWMutex
}

func newMemoryCacheEntry() *memoryCacheEntry {
	return &memoryCacheEntry{}
}

func (entry *memoryCacheEntry) Value() (value interface{}, err error) {
	// 如果过期了
	if !entry.isExpired() {
		entry.keep()
	}
	return entry.value, nil
}

func (entry *memoryCacheEntry) isExpired() bool {
	return !entry.policy.AbsoluteExpiration.IsZero() && time.Now().After(entry.policy.AbsoluteExpiration)
}

func (entry *memoryCacheEntry) keep() {
	if entry.policy.SlidingExpiration > 0 {
		entry.policy.AbsoluteExpiration = time.Now().Add(entry.policy.SlidingExpiration)
	}
}
func (entry *memoryCacheEntry) hasExpiration() bool {
	return entry.policy.SlidingExpiration > 0 || entry.policy.AbsoluteExpiration != MaxTimeValue
}

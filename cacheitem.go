package gocache

import (
	"sync"
	"time"
)

type CacheItem struct {
	sync.RWMutex
	key               string
	value             interface{}
	absExp            time.Time
	slidingExpiration time.Duration
	UpdateCallback    CacheEntryUpdateCallback
	RemovedCallback   CacheEntryRemovedCallback
	CreateCallback    CacheEntryCreateCallback
}

func NewCacheItem(key string, value interface{}) *CacheItem {
	return &CacheItem{
		key:    key,
		value:  value,
		absExp: MaxTimeValue,
	}
}

func NewCacheItemWithSlidingExpiration(key string, value interface{}, expiration time.Duration) *CacheItem {
	return &CacheItem{
		key:               key,
		value:             value,
		absExp:            time.Now().Add(expiration),
		slidingExpiration: expiration,
	}
}

func NewCacheItemWithAbsoluteExpiration(key string, value interface{}, expiration time.Time) *CacheItem {
	return &CacheItem{
		key:    key,
		value:  value,
		absExp: expiration,
	}
}

// 是否存在过期日期设置
func (item *CacheItem) HasExpiration() bool {
	return !item.absExp.Equal(MaxTimeValue) && item.absExp.Before(MaxTimeValue)
}

// 是否已过期
func (item *CacheItem) InExpires() bool {
	return item.absExp.Before(time.Now())
}

// 更新过期日期
func (item *CacheItem) UpdateSlidingExpiration(expiration time.Duration) {
	item.absExp = time.Now().Add(expiration)
	item.slidingExpiration = expiration
}

//更新绝对过期时间
func (item *CacheItem) UpdateAbsoluteExpiration(expiration time.Time) {
	item.absExp = expiration
}

func (item *CacheItem) KeepLive() {
	if item.slidingExpiration > 0 {
		item.absExp = time.Now().Add(item.slidingExpiration)
	}
}

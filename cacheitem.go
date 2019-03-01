package gocache

import (
	"sync"
	"time"
)

type CacheItem struct {
	sync.RWMutex
	Key               string
	Value             interface{}
	absExp            time.Time
	slidingExpiration time.Duration
	RemovedCallback   CacheEntryRemovedCallback
	CreateCallback    CacheEntryCreateCallback
}

func NewCacheItem(key string, value interface{}) *CacheItem {
	return &CacheItem{
		Key:    key,
		Value:  value,
		absExp: MaxTimeValue,
	}
}

func NewCacheItemWithSlidingExpiration(key string, value interface{}, expiration time.Duration) *CacheItem {
	return &CacheItem{
		Key:               key,
		Value:             value,
		slidingExpiration: expiration,
	}
}

func NewCacheItemWithAbsoluteExpiration(key string, value interface{}, expiration time.Time) *CacheItem {
	return &CacheItem{
		Key:    key,
		Value:  value,
		absExp: expiration,
	}
}

// 是否存在过期日期设置
func (item *CacheItem) HasExpiration() bool {
	return (!item.absExp.Equal(MaxTimeValue) && item.absExp.Before(MaxTimeValue)) || item.slidingExpiration > 0
}

// 是否已过期
func (item *CacheItem) InExpires() bool {
	return item.absExp.Before(time.Now())
}

// 更新过期日期
func (item *CacheItem) updateSlidingExpiration(expiration time.Duration) {
	item.absExp = time.Now().Add(expiration)
	item.slidingExpiration = expiration
}

//更新绝对过期时间
func (item *CacheItem) updateAbsoluteExpiration(expiration time.Time) {
	item.absExp = expiration
}

func (item *CacheItem) keepLive() {
	if item.slidingExpiration > 0 {
		item.absExp = time.Now().Add(item.slidingExpiration)
	}
}

func (item *CacheItem) callRemovedCallback(reason CacheEntryRemovedReason) {
	if item.RemovedCallback != nil {
		item.RemovedCallback(item.Key, item.Value, reason)
	}
}

func (item *CacheItem) reset() *CacheItem {
	if item == nil {
		return item
	}
	item.Key = ""
	item.Value = nil
	item.CreateCallback = nil
	item.RemovedCallback = nil
	item.absExp = MaxTimeValue
	item.slidingExpiration = time.Duration(0)
	return item
}

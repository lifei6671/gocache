package gocache

import (
	"time"
)

type memoryCacheEntry struct {
	Key     string
	value   interface{}
	created time.Time
	policy  *CacheItemPolicy
}

func newMemoryCacheEntry() *memoryCacheEntry {
	return &memoryCacheEntry{}
}

func (entry *memoryCacheEntry) Value() (value interface{}, err error) {
	// 如果过期了
	if !entry.policy.AbsoluteExpiration.IsZero() && time.Now().Before(entry.policy.AbsoluteExpiration) {
		//必须设置了滑动过期时间才能通过回调方法创建缓存值
		if entry.policy.CreateCallback != nil && entry.policy.SlidingExpiration > 0 {
			if v, err := entry.policy.CreateCallback(entry.Key); err != nil {
				return nil, err
			} else {
				entry.value = v
			}
		}
	}
	//更新最后过期时间
	entry.policy.AbsoluteExpiration = time.Now().Add(entry.policy.SlidingExpiration)
	return entry.value, nil
}

func (entry *memoryCacheEntry) isExpired() bool {
	return !entry.policy.AbsoluteExpiration.IsZero() && time.Now().After(entry.policy.AbsoluteExpiration)
}

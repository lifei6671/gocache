package gocache

import "time"

var defaultMemoryCache = NewMemoryCache(time.Duration(1))

func Add(key string, value interface{}) {
	store := defaultMemoryCache.getStore(key)
	item := NewCacheItem(key, value)
	store.Add(item)
}

func AddWithCacheItem(item CacheItem) {
	store := defaultMemoryCache.getStore(item.Key)
	store.Add(&item)
}

// 增加缓存并设置策略
func AddWithSlidingExpiration(key string, value interface{}, expiration time.Duration) {
	store := defaultMemoryCache.getStore(key)
	store.AddWithSlidingExpiration(key, value, expiration)
}

func AddWithAbsoluteExpiration(key string, value interface{}, expiration time.Time) {
	store := defaultMemoryCache.getStore(key)
	store.AddWithAbsoluteExpiration(key, value, expiration)
}

func Get(key string) (value interface{}, ok bool) {
	store := defaultMemoryCache.getStore(key)
	return store.GetValue(key)
}

// ContainsKey 判断缓存中是否存在指定键
func ContainsKey(key string) bool {
	store := defaultMemoryCache.getStore(key)
	return store.ContainsKey(key)
}

func Count() int {
	c := 0
	for _, store := range defaultMemoryCache.stores {
		c += store.Count()
	}
	return c
}

package gocache

//定义对某个方法的引用，在从缓存中移除某个缓存项后将调用该方法
type CacheEntryRemovedCallback func(key string, value interface{}, reason CacheEntryRemovedReason)

// 定义对某个方法的引用，在即将从缓存中移除某个缓存项时将调用该方法
type CacheEntryUpdateCallback func(key string, value interface{}, reason CacheEntryRemovedReason)

// 当缓存失效后，可以指定重新生成缓存的回调方法
type CacheEntryCreateCallback func(key string, oldValue interface{}) (value interface{}, err error)

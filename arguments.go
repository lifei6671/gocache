package gocache

type CacheEntryRemovedArguments struct {
	//获取已从缓存中移除的缓存项的实例
	CacheItem CacheItem
	//获取一个值，该值指示移除某个缓存项的原因
	RemovedReason CacheEntryRemovedReason
}

type CacheEntryUpdateArguments struct {
	//从缓存中移除某个缓存项的原因
	RemovedReason CacheEntryRemovedReason
	//获取或设置用于更新缓存对象的 CacheItem 项的值
	UpdatedCacheItem CacheItem
}

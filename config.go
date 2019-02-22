package gocache

type CacheEntryUpdateArguments struct {
	//将要移除的缓存项的唯一标识符
	Key string
	//包含某个缓存项的缓存区域的名称
	RegionName string
	//从缓存中移除某个缓存项的原因
	RemovedReason CacheEntryRemovedReason
	//包含一个将要移除的缓存项
	Source interface{}
	//获取或设置用于更新缓存对象的 CacheItem 项的值
	UpdatedCacheItem *CacheItem
	//获取或设置已更新的 CacheItem 项的缓存逐出或过期策略
	UpdatedCacheItemPolicy
}

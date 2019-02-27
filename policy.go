package gocache

import (
	"time"
)

//https://docs.microsoft.com/zh-cn/dotnet/api/system.runtime.caching.cacheitempolicy?view=netframework-4.7.2
type CacheItemPolicy struct {
	//该值指示是否应在指定的时间点逐出缓存项
	AbsoluteExpiration time.Time
	//获取或设置用于确定是否逐出某个缓存项的优先级别设置
	Priority CacheItemPriority

	//缓存中移除某个项后将调用
	RemovedCallback CacheEntryRemovedCallback

	//获取或设置一个值，该值指示如果某个缓存项在给定时段内未被访问，是否应被逐出
	SlidingExpiration time.Duration

	// 缓存中移除某个缓存项之前将调用
	UpdateCallback CacheEntryUpdateCallback

	//当缓存失效后调用该方法重新生成缓存
	CreateCallback CacheEntryCreateCallback
}

func NewCacheItemPolicy() CacheItemPolicy {
	return CacheItemPolicy{
		AbsoluteExpiration: MaxTimeValue,
	}
}

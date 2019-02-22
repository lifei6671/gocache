package gocache

import "time"

//https://docs.microsoft.com/zh-cn/dotnet/api/system.runtime.caching.cacheitempolicy?view=netframework-4.7.2
type CacheItemPolicy struct {
	//该值指示是否应在指定的时间点逐出缓存项
	AbsoluteExpiration time.Time
}

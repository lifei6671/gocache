package gocache

import "time"

var MaxTimeValue = time.Date(9999, 12, 31, 23, 59, 59, 9999999, time.Local)

//指定已移除或将要移除某个缓存项的原因
type CacheEntryRemovedReason int

const (
	// 主动移除
	CacheEntryRemovedReasonRemoved CacheEntryRemovedReason = iota
	//某个缓存项由于已过期而被移除。 过期可基于绝对过期时间或可调过期时间
	CacheEntryRemovedReasonExpired

	//某个缓存项由于释放缓存中的内存的原因而被移除。
	// 当某个缓存实例将超出特定于缓存的内存限制或某个进程或缓存实例将超出整个计算机范围的内存限制时，会发生这种情况
	CacheEntryRemovedReasonEvicted

	//某个缓存项由于相关依赖项（如一个文件或其他缓存项）触发了其逐出操作而被移除
	CacheEntryRemovedReasonChangeMonitorChanged

	//某个缓存项由于特定缓存实现定义的原因而被逐出
	CacheEntryRemovedReasonCacheSpecificEviction
)

// 一个时段，必须在此时段内访问某个缓存项，否则将从内存中逐出该缓存项
type SlidingExpiration time.Duration

// 基于持续时间或定义的时间窗口过期也被称为可调过期。 通常情况下，逐出基于可调过期的项的缓存实现将删除在指定的时间段未被访问的项。
// 插入到缓存与某个缓存项NoSlidingExpiration字段值设置为过期值应该永远不会收回由于非活动的滑动时间窗口中。
// 但是，如果具有绝对到期，或某些其他逐出事件发生时，此类更改监视器或内存压力，可以逐出缓存项。
type NoSlidingExpiration time.Duration

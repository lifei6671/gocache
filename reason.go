package gocache

const (
	//主动删除
	Removed CacheEntryRemovedReason = iota
	//某个缓存项由于已过期而被移除。 过期可基于绝对过期时间或可调过期时间
	Expired
	//某个缓存项由于释放缓存中的内存的原因而被移除。
	// 当某个缓存实例将超出特定于缓存的内存限制或某个进程或缓存实例将超出整个计算机范围的内存限制时，会发生这种情况
	Evicted
	//某个缓存项由于相关依赖项（如一个文件或其他缓存项）触发了其逐出操作而被移除
	ChangeMonitorChanged
	//某个缓存项由于特定缓存实现定义的原因而被逐出。
	CacheSpecificEviction
)

type CacheEntryRemovedReason int

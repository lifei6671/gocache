package gocache

const (
	// 指示移除缓存项没有优先级
	CacheItemPriorityDefault CacheItemPriority = iota

	//指示绝不应从缓存中移除某个缓存项
	CacheItemPriorityNotRemovable
)

//指定用于确定是否逐出某个缓存项的优先级别设置
type CacheItemPriority int

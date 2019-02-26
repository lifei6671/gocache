package gocache

var SHARD_COUNT = 32

type MemoryCache struct {
	stores     []*MemoryCacheStore
	storeCount int
	closed     bool
}

func NewMemoryCache() *MemoryCache {

	cache := &MemoryCache{
		stores:     make([]*MemoryCacheStore, SHARD_COUNT),
		storeCount: SHARD_COUNT,
	}

	for i := 0; i < SHARD_COUNT; i++ {
		cache.stores[i] = &MemoryCacheStore{entries: make(map[string]*memoryCacheEntry, 10000)}
	}

	return cache
}

func (memory *MemoryCache) Add(key string, value interface{}) {
	store := memory.getStore(key)
	store.Add(key, value)
}

// 增加缓存并设置策略
func (memory *MemoryCache) AddWithPolicy(key string, value interface{}, policy CacheItemPolicy) {
	store := memory.getStore(key)
	store.AddWithPolicy(key, value, &policy)
}

func (memory *MemoryCache) Get(key string) (value interface{}, ok bool) {
	store := memory.getStore(key)
	return store.Get(key)
}

// ContainsKey 判断缓存中是否存在指定键
func (memory *MemoryCache) ContainsKey(key string) bool {
	store := memory.getStore(key)
	return store.ContainsKey(key)
}

func (memory *MemoryCache) Count() int {
	c := 0
	for _, store := range memory.stores {
		c += store.Count()
	}
	return c
}

func (memory *MemoryCache) getStore(key string) *MemoryCacheStore {
	idx := int(hashKey(key) % uint32(memory.storeCount))

	return memory.stores[idx]
}

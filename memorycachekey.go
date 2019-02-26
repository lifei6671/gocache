package gocache

type MemoryCacheKey struct {
	Hash int64
	Key  string
}

func NewMemoryCacheKey(key string) *MemoryCacheKey {
	cacheKey := &MemoryCacheKey{
		Key: key,
	}

	return cacheKey
}

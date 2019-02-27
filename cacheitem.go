package gocache

import "fmt"

type CacheItem struct {
	Key   string
	Value interface{}
}

func (item CacheItem) String() string {
	return fmt.Sprintf("Key:%s; Value:%+v", item.Key, item.Value)
}

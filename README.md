# gocache

## 简介

gocache 使用分片 + 锁实现进程内缓存。

缓存根据键的hash值分片到多个储存空间，每个储存空间单独加锁。

1. 实现了滑动过期和绝对过期
2. 实现了滑动过期指定创建函数再次生成缓存
3. 实现缓存过期删除时的回调事件

## 使用

### 使用内置缓存

```go
package main

import (
	"github.com/lifei6671/gocache"
	"log"
	"time"
)

func main() {
	
	gocache.AddWithCacheItem(gocache.CacheItem{Key: "cache", Value: "cache-1", RemovedCallback: func(key string, value interface{}, reason gocache.CacheEntryRemovedReason) {
		log.Println(key, value, reason)
	}})
	gocache.Add("cache", "cache-2")

	log.Println(gocache.Count())
	log.Println(gocache.Get("cache"))

}
```

### 使用自定义缓存

```go
package main

import (
	"github.com/lifei6671/gocache"
	"log"
	"time"
)

func main() {
	m := gocache.NewMemoryCache(time.Second * 2)

	m.AddWithCacheItem(gocache.CacheItem{Key: "cache", Value: "cache-1", RemovedCallback: func(key string, value interface{}, reason gocache.CacheEntryRemovedReason) {
		log.Println(key, value, reason)
	}})
	m.Add("cache", "cache-2")

	log.Println(m.Count())
	log.Println(m.Get("cache"))

}

```

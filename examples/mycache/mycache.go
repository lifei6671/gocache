package main

import (
	"github.com/lifei6671/gocache"
	"log"
	"time"
)

// Keys & values in cache2go can be of arbitrary types, e.g. a struct.
type myStruct struct {
	text     string
	moreData []byte
}

func main() {
	m := gocache.NewMemoryCache()

	val := myStruct{"This is a test!", []byte{}}

	m.Add("someKey", val)

	v, ok := m.Get("someKey")

	if !ok {
		log.Fatal("缓存不存在")
	} else if my, ok := v.(myStruct); !ok {
		log.Fatal("缓存获取失败")
	} else {
		log.Println("缓存获取成功", my.text)
	}

	policy := gocache.CacheItemPolicy{
		SlidingExpiration: time.Second * 10,
		CreateCallback: func(key string) (value interface{}, err error) {
			return "aaaaaaaaaaaa", nil
		},
		RemovedCallback: func(arguments gocache.CacheEntryRemovedArguments) {
			log.Println(arguments.RemovedReason)
		},
	}
	m.AddWithPolicy("someKey", 10, policy)

	v, ok = m.Get("someKey")

	if !ok {
		log.Fatal("缓存不存在")
	} else {
		log.Println("获取成功", v)
	}
	time.AfterFunc(time.Second*11, func() {
		v, ok = m.Get("someKey")

		if !ok {
			log.Fatal("缓存不存在")
		} else {
			log.Println("获取成功", v)
		}
	})

	select {}
}

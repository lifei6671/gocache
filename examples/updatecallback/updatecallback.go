package main

import (
	"fmt"
	"github.com/lifei6671/gocache"
	"log"
	"time"
)

type myStruct struct {
	text     string
	moreData []byte
}

func (m *myStruct) String() string {
	return fmt.Sprintf("text:%s; moreData:%s", m.text, string(m.moreData))
}

func main() {
	m := gocache.NewMemoryCache()

	go func() {
		for i := 0; i < 250; i++ {
			m.AddWithPolicy(fmt.Sprintf("sso-%d", i), &myStruct{text: fmt.Sprintf("sso-%d", i), moreData: []byte(fmt.Sprintf("ssox系统%d", i))}, gocache.CacheItemPolicy{
				UpdateCallback: func(arguments gocache.CacheEntryUpdateArguments) {
					log.Println("缓存已更新 ->", arguments.RemovedReason, arguments.UpdatedCacheItem)
				},
				RemovedCallback: func(arguments gocache.CacheEntryRemovedArguments) {
					log.Println("缓存已删除 ->", arguments.RemovedReason, arguments.CacheItem)
				},
				SlidingExpiration: time.Second * time.Duration(i+1),
			})
		}
	}()
	go func() {
		for i := 250; i < 500; i++ {
			m.AddWithPolicy(fmt.Sprintf("sso-%d", i), &myStruct{text: fmt.Sprintf("sso-%d", i), moreData: []byte(fmt.Sprintf("ssox系统%d", i))}, gocache.CacheItemPolicy{
				UpdateCallback: func(arguments gocache.CacheEntryUpdateArguments) {
					log.Println("缓存已更新 ->", arguments.RemovedReason, arguments.UpdatedCacheItem)
				},
				RemovedCallback: func(arguments gocache.CacheEntryRemovedArguments) {
					log.Println("缓存已删除 ->", arguments.RemovedReason, arguments.CacheItem)
				},
				SlidingExpiration: time.Second * time.Duration(i+1),
			})
		}
	}()
	go func() {
		for i := 500; i < 750; i++ {
			m.AddWithPolicy(fmt.Sprintf("sso-%d", i), &myStruct{text: fmt.Sprintf("sso-%d", i), moreData: []byte(fmt.Sprintf("ssox系统%d", i))}, gocache.CacheItemPolicy{
				UpdateCallback: func(arguments gocache.CacheEntryUpdateArguments) {
					log.Println("缓存已更新 ->", arguments.RemovedReason, arguments.UpdatedCacheItem)
				},
				RemovedCallback: func(arguments gocache.CacheEntryRemovedArguments) {
					log.Println("缓存已删除 ->", arguments.RemovedReason, arguments.CacheItem)
				},
				SlidingExpiration: time.Second * time.Duration(i+1),
			})
		}
	}()
	go func() {
		for i := 750; i < 1000; i++ {
			m.AddWithPolicy(fmt.Sprintf("sso-%d", i), &myStruct{text: fmt.Sprintf("sso-%d", i), moreData: []byte(fmt.Sprintf("ssox系统%d", i))}, gocache.CacheItemPolicy{
				UpdateCallback: func(arguments gocache.CacheEntryUpdateArguments) {
					log.Println("缓存已更新 ->", arguments.RemovedReason, arguments.UpdatedCacheItem)
				},
				RemovedCallback: func(arguments gocache.CacheEntryRemovedArguments) {
					log.Println("缓存已删除 ->", arguments.RemovedReason, arguments.CacheItem)
				},
				SlidingExpiration: time.Second * time.Duration(i+1),
			})
		}
	}()

	m.Add("rbac", &myStruct{text: "rbac", moreData: []byte("rbac系统")})
	m.AddWithPolicy("sso", &myStruct{text: "sso", moreData: []byte("ssox系统")}, gocache.CacheItemPolicy{
		UpdateCallback: func(arguments gocache.CacheEntryUpdateArguments) {
			log.Println("缓存已更新 ->", arguments.RemovedReason, arguments.UpdatedCacheItem)
		},
		RemovedCallback: func(arguments gocache.CacheEntryRemovedArguments) {
			log.Println("缓存已删除 ->", arguments.RemovedReason, arguments.CacheItem)
		},
		SlidingExpiration: time.Second * 10,
	})

	m.AddWithPolicy("mc", &myStruct{text: "mc", moreData: []byte("mc系统")}, gocache.CacheItemPolicy{
		UpdateCallback: func(arguments gocache.CacheEntryUpdateArguments) {
			log.Println("缓存已更新 ->", arguments.RemovedReason, arguments.UpdatedCacheItem)
		},
		RemovedCallback: func(arguments gocache.CacheEntryRemovedArguments) {
			log.Println("缓存已删除 ->", arguments.RemovedReason, arguments.CacheItem)
		},
		AbsoluteExpiration: time.Now().Add(time.Second * 10),
	})

	timer := time.NewTimer(time.Second * 2)

	for {
		select {
		case <-timer.C:
			m.ContainsKey("sso-10")

			log.Println(m.Get("sso"))
			log.Println(m.Get("sso"))
			log.Println(m.Get("rbac"))
			log.Println(m.ContainsKey("sso-20"))
			if v, ok := m.Get("mc"); ok {
				log.Println(v, ok)
				timer.Reset(time.Second * 2)
			} else {
				timer.Reset(time.Second * 12)
			}
		}
	}

}

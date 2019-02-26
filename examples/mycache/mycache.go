package main

import (
	"github.com/lifei6671/gocache"
	"log"
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
		log.Println(my.text)
	}
}

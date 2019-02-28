package main

import (
	"github.com/lifei6671/gocache"
	"log"
	"time"
)

func main() {
	m := gocache.NewMemoryCache(time.Second * 2)

	m.Add("cache", "cache-1")
	m.Add("cache", "cache-2")

	log.Println(m.Count())
	log.Println(m.Get("cache"))

}

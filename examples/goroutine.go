package main

import (
	"fmt"
	"github.com/lifei6671/gocache"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type User struct {
	Name string
}

func main() {
	//这里实现了远程获取pprof数据的接口
	go func() {
		log.Println(http.ListenAndServe("localhost:7777", nil))
	}()
	m := gocache.NewMemoryCache(time.Second * 2)

	m.Add("tools-all", User{Name: "tools-all"})
	for i := 1; i < 200; i++ {
		go func(i int) {
			for j := i * 10; j < i*100; j++ {
				m.Add(fmt.Sprintf("tools-%d", j), User{Name: fmt.Sprintf("user-%d", j)})
			}
		}(i)
	}

	timer := time.NewTimer(time.Second * 2)

	for {
		select {
		case <-timer.C:
			m.ContainsKey("tools-all")
			log.Println(m.Count())
			log.Println(m.Get("tools-100"))
			timer.Reset(time.Second * 2)
		}
	}
}

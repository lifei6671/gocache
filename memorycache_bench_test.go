package gocache

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkMemoryCache_Add(b *testing.B) {
	m := NewMemoryCache(time.Second * 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("cache-%d", i)

		m.Add(k, "开源软件没你想象中那么安全，Java 开发者尤其要警惕-%d")
	}
}

func BenchmarkMemoryCache_AddWithSlidingExpiration(b *testing.B) {
	m := NewMemoryCache(time.Second * 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("cache-%d", i)

		m.AddWithSlidingExpiration(k, "开源软件没你想象中那么安全，Java 开发者尤其要警惕-%d", time.Second*10)
	}
}

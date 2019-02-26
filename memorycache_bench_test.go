package gocache

import (
	"math/rand"
	"strconv"
	"testing"
)

func nrand(n int) []int {
	i := make([]int, n)
	for ind := range i {
		i[ind] = rand.Int()
	}
	return i
}

func BenchmarkMemoryCache_Add(b *testing.B) {
	m := NewMemoryCache()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k := strconv.Itoa(i)
		m.Add(k, Animal{k})
	}
}

func BenchmarkMemoryCache_Get(b *testing.B) {
	m := NewMemoryCache()

	for i := 0; i < 1000; i++ {
		k := strconv.Itoa(i)
		m.Add(k, Animal{k})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := strconv.Itoa(i)
		m.Get(k)
	}
}

package cache

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkMemoryCache_Store(b *testing.B) {
	counts := []int{8192, 4096, 2048, 1024, 512, 256}

	for _, chunkCounts := range counts {
		b.Run(fmt.Sprintf("chunk-%d", chunkCounts), func(b *testing.B) {
			memoryCache := New(30*time.Minute, 1*time.Hour, chunkCounts)
			memoryCache.Run()

			b.Cleanup(func() {
				memoryCache.Stop()
			})

			key := make([]byte, 64)
			val := bytes.Repeat([]byte{'$'}, 4096)

			b.RunParallel(func(pb *testing.PB) {
				r := rand.New(rand.NewSource(time.Now().Unix()))
				b.ReportAllocs()
				for pb.Next() {
					r.Read(key)
					memoryCache.Store(string(key), val)
				}
			})
		})
	}
}

func BenchmarkMemoryCache_Get(b *testing.B) {
	counts := []int{8192, 4096, 2048, 1024, 512, 256}
	for _, chunkCounts := range counts {
		b.Run(fmt.Sprintf("chunk-%d", chunkCounts), func(b *testing.B) {
			memoryCache := New(30*time.Minute, 1*time.Hour, chunkCounts)
			memoryCache.Run()

			b.Cleanup(func() {
				memoryCache.Stop()
			})

			val := bytes.Repeat([]byte{'$'}, 4096)

			for i := 0; i < b.N; i++ {
				memoryCache.Store(strconv.Itoa(i), val)
			}
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				r := rand.New(rand.NewSource(time.Now().Unix()))
				b.ReportAllocs()

				for pb.Next() {
					memoryCache.Get(strconv.Itoa(r.Intn(b.N + 1)))
				}
			})
		})
	}
}

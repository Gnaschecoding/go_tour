package cache

import (
	"5_cache/lru"
	"fmt"
	"github.com/matryer/is"
	"log"
	"sync"
	"testing"
)

// 构造通用的 TourCache 用于基准测试
func buildBenchmarkCache() *TourCache {
	db := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
	}
	getter := GetFunc(
		func(key string) interface {
		} {
			if v, ok := db[key]; ok {
				return v
			}
			return nil
		})
	return NewTourCache(getter, lru.New(1000, nil))
}

func BenchmarkTourCacheSet(b *testing.B) {
	cache := buildBenchmarkCache()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i))
	}
}

func BenchmarkTourCacheGet(b *testing.B) {
	cache := buildBenchmarkCache()
	// 预先放入数据
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cache.Get(fmt.Sprintf("key%d", i%1000))
	}
}

func BenchmarkTourCacheSetParallel(b *testing.B) {
	cache := buildBenchmarkCache()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i))
			i++
		}
	})
}

func BenchmarkTourCacheGetParallel(b *testing.B) {
	cache := buildBenchmarkCache()
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i))
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = cache.Get(fmt.Sprintf("key%d", i%1000))
			i++
		}
	})
}

func TestTourCacheGet(t *testing.T) {
	db := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
	}
	getter := GetFunc(func(key string) interface{} {
		if v, ok := db[key]; ok {
			log.Println("[From DB] find key", key)
			return v
		}
		return nil
	})

	tourCache := NewTourCache(getter, lru.New(1000, nil))
	is := is.New(t)

	var wg sync.WaitGroup
	for k, v := range db {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			is.Equal(tourCache.Get(k), v)
			is.Equal(tourCache.Get(k), v)
		}(k, v)
	}
	wg.Wait()

	is.Equal(tourCache.Get("unknown"), nil)
	is.Equal(tourCache.Get("unknown"), nil)

	is.Equal(tourCache.Stat().NGet, 10)
	is.Equal(tourCache.Stat().NHit, 4)

}

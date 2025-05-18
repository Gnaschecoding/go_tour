package cache

import (
	"5_cache/lru"
	"github.com/matryer/is"
	"log"
	"sync"
	"testing"
)

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

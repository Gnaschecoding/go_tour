package main

import (
	"github.com/allegro/bigcache"
	"log"
	"time"
)

func main() {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		log.Println(err)
		return
	}

	entry, err := cache.Get("my-unique-key")
	if err != nil {
		log.Println(err)
		return
	}

	if entry == nil {
		// 从缓存中没有获取到，则从数据源取（一般是数据库），然后设置到缓存
		entry = []byte("value") // 实际从数据库获取
		cache.Set("my-unique-key", entry)
	}

	log.Println(string(entry))
}

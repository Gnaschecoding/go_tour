package cache

import (
	"log"
	"sync"
)

// Cache 缓存接口
type Cache interface {
	Set(key string, value interface{}) //设置/添加一个缓存，如果 key 存在，用新值覆盖旧值；
	Get(key string) interface{}        //	通过 key 获取一个缓存值；
	Del(key string)                    //	通过 key 删除一个缓存值；
	DelOldest()                        //删除最“无用”的一个缓存值；
	Len() int
}

// DefaultMaxBytes 默认允许占用的最大内存
const DefaultMaxBytes = 1 << 29

// safeCache 并发安全缓存
type SafeCache struct {
	m     sync.RWMutex
	cache Cache

	nget, nhit int
}

func NewSafeCache(cache Cache) *SafeCache {
	return &SafeCache{
		cache: cache,
	}
}

func (sc *SafeCache) Set(key string, value interface{}) {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.cache.Set(key, value)
}

func (sc *SafeCache) Get(key string) interface{} {
	sc.m.RLock()
	defer sc.m.RUnlock()
	sc.nget++

	if sc.cache == nil {
		return nil
	}
	v := sc.cache.Get(key)
	if v != nil {
		log.Println("[TourCache] hit")
		sc.nhit++
	}
	return v
}

func (sc *SafeCache) Stat() *Stat {
	sc.m.RLock()
	defer sc.m.RUnlock()
	return &Stat{
		NGet: sc.nget,
		NHit: sc.nhit,
	}
}

type Stat struct {
	NHit, NGet int
}

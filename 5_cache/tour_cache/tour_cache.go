package cache

import cache "5_cache"

type Getter interface {
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}

type TourCache struct {
	mainCache *cache.SafeCache
	getter    Getter
}

func NewTourCache(getter Getter, cache_ cache.Cache) *TourCache {
	return &TourCache{
		mainCache: cache.NewSafeCache(cache_),
		getter:    getter,
	}
}

func (t *TourCache) Get(key string) interface{} {
	val := t.mainCache.Get(key)
	if val != nil {
		return val
	}

	if t.getter != nil {
		val = t.getter.Get(key)
		if val != nil {
			t.mainCache.Set(key, val)
			return val
		}
	}
	return nil

}

func (t *TourCache) Stat() *cache.Stat {
	return t.mainCache.Stat()
}

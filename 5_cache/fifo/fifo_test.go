package fifo

import (
	"github.com/matryer/is"
	"testing"
)

func TestSetGet(t *testing.T) {
	is := is.New(t)

	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)

	v := cache.Get("k1")
	is.Equal(v, 1)

	cache.Del("k1")
	is.Equal(0, cache.Len()) // expect to be the same
}

func TestOnEvicted(t *testing.T) {
	is := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := New(16, onEvicted)

	cache.Set("k1", 1) //其实占据了 10个字节空间  k1长度占2，1是int64占用8个字节
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k1")
	cache.Set("k4", 4)

	expected := []string{"k1", "k2"}

	is.Equal(expected, keys)
	is.Equal(2, cache.Len())
}

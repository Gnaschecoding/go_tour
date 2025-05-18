package lru

import (
	cache "5_cache"
	"container/list"
)

// fifo 是一个 FIFO cache。它不是并发安全的。
type lru struct {
	// 缓存最大的容量，单位字节；
	// groupcache 使用的是最大存放 entry 个数
	maxBytes int
	// 当一个 entry 从缓存中移除是调用该回调函数，默认为 nil
	// groupcache 中的 key 是任意的可比较类型；value 是 interface{}
	onEvicted func(key string, value interface{})

	// 已使用的字节数，只包括值，key 不算
	usedBytes int

	ll    *list.List
	cache map[string]*list.Element
}

// New 创建一个新的 Cache，如果 maxBytes 是 0，表示没有容量限制
func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &lru{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return len(e.key) + cache.CalcLen(e.value)
}

// Set 往 Cache 尾部增加一个元素（如果已经存在，则移到尾部，并修改值）
func (l *lru) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToFront(e)
		en := e.Value.(*entry)
		l.usedBytes = l.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
		en.value = value
		for l.maxBytes > 0 && l.usedBytes > l.maxBytes {
			l.DelOldest()
		}
		return
	}

	en := &entry{key, value}
	e := l.ll.PushFront(en)
	l.cache[key] = e

	l.usedBytes += en.Len()
	for l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

// Get 从 cache 中获取 key 对应的值，nil 表示 key 不存在
func (l *lru) Get(key string) interface{} {
	if l.cache == nil {
		return nil
	}
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToFront(e)
		return e.Value.(*entry).value
	}
	return nil
}

// Del 从 cache 中删除 key 对应的记录
func (l *lru) Del(key string) {
	if l.cache == nil || len(l.cache) == 0 {
		return
	}
	if e, ok := l.cache[key]; ok {
		l.removeElement(e)
	}
	return
}

// DelOldest 从 cache 中删除最旧的记录'
func (l *lru) DelOldest() {
	if l.cache == nil || len(l.cache) == 0 {
		return
	}
	l.removeElement(l.ll.Back())
}
func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	l.ll.Remove(e)
	en := e.Value.(*entry)
	l.usedBytes -= en.Len()
	delete(l.cache, en.key)

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

// Len 返回当前 cache 中的记录数
func (l *lru) Len() int {
	return l.ll.Len()
}

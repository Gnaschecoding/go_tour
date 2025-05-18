package fast

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const maxEntrySize = 256

func BenchmarkTourFastCacheSetParallel(b *testing.B) {
	cache := NewFastCache(b.N, maxEntrySize, nil)
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set(parallelKey(id, counter), value())
			counter++
		}
	})
}

func value() []byte {
	return make([]byte, 100)
}

func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
}

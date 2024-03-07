package cache

import (
	"sync"
	"time"
)

// A chunk is a partition of a wider memory cache, which is thread-safe to access.
// Rather than store data in a single, massive map, spread it out for quicker retrieval.
type chunk struct {
	mut    sync.RWMutex
	store  map[uint64][]byte
	expire map[uint64]time.Time
}

func (ch *chunk) Get(key uint64) ([]byte, time.Time, bool) {
	ch.mut.RLock()
	defer ch.mut.RUnlock()

	val, ok := ch.store[key]
	expire, _ := ch.expire[key]
	return val, expire, ok
}

func (ch *chunk) Store(key uint64, val []byte, expiration time.Time) {
	// Opting to overwrite existing cache values and ignore what was previously there.
	// It might be useful to pop the existing entry instead, but that's for the future
	ch.mut.Lock()
	defer ch.mut.Unlock()

	ch.store[key] = val
	ch.expire[key] = expiration
}

func (ch *chunk) Delete(key uint64) {
	ch.mut.Lock()
	defer ch.mut.Unlock()

	delete(ch.store, key)
	delete(ch.expire, key)
}

func (ch *chunk) GC() {
	ch.mut.RLock()

	currentTime := time.Now()
	evict := make([]uint64, 0)

	for k, v := range ch.expire {
		if v.Before(currentTime) {
			evict = append(evict, k)
		}
	}
	ch.mut.RUnlock()

	for _, key := range evict {
		ch.Delete(key)
	}
}

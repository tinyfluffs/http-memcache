package cache

import (
	"hash/maphash"
	"sync"
	"time"
)

type MemoryCache struct {
	expiration time.Duration
	gcInterval time.Duration
	done       chan struct{}
	mask       uint64
	chunk      []*chunk
	seed       maphash.Seed
	gcTicker   *time.Ticker
}

func New(expiration time.Duration, gcInterval time.Duration, chunkCount int) *MemoryCache {
	mc := &MemoryCache{
		expiration: expiration,
		gcInterval: gcInterval,
		done:       make(chan struct{}),
		mask:       uint64(chunkCount - 1),
		chunk:      make([]*chunk, chunkCount),
		seed:       maphash.MakeSeed(),
		gcTicker:   time.NewTicker(gcInterval),
	}

	for i := 0; i < chunkCount; i++ {
		mc.chunk[i] = &chunk{
			mut:    sync.RWMutex{},
			store:  make(map[uint64][]byte),
			expire: make(map[uint64]time.Time),
		}
	}

	return mc
}

func (mc *MemoryCache) Run() {
	go func() {

		for {
			select {
			case <-mc.done:
				mc.gcTicker.Stop()
				return
			case <-mc.gcTicker.C:
				mc.GC()
			}
		}
	}()
}

func (mc *MemoryCache) Stop() {
	close(mc.done)
}

func (mc *MemoryCache) Get(k string) ([]byte, time.Time, bool) {
	key := mc.hash(k)
	chk := mc.chunkForKey(key)
	val, expiry, ok := chk.Get(key)
	if expiry.Before(time.Now()) {
		go chk.Delete(key) // Cleanup now, why wait for GC to hit
		return val, expiry, false
	}
	return val, expiry, ok
}

func (mc *MemoryCache) Store(k string, v []byte) {
	key := mc.hash(k)
	chk := mc.chunkForKey(key)
	chk.Store(key, v, time.Now().Add(mc.expiration))
}

func (mc *MemoryCache) GC() {
	for _, c := range mc.chunk {
		c.GC()
	}
}

func (mc *MemoryCache) chunkForKey(keyHash uint64) *chunk {
	return mc.chunk[keyHash&mc.mask]
}

func (mc *MemoryCache) hash(key string) uint64 {
	return maphash.String(mc.seed, key)
}

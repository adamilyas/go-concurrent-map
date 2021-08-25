package cmap

import (
	"sync"
)

const shardCount = 32

type ConcurrentMap struct {
	shards     []*ConcurrentMapShard
	shardCount int
}

type ConcurrentMapShard struct {
	items map[string]interface{}
	sync.RWMutex
}

func New(shardCount int) *ConcurrentMap {
	m := make([]*ConcurrentMapShard, shardCount)

	for i := 0; i < shardCount; i++ {
		m[i] = &ConcurrentMapShard{items: map[string]interface{}{}}
	}

	return &ConcurrentMap{shards: m, shardCount: shardCount}
}

func (m *ConcurrentMap) Count() int {
	count := 0

	for i := 0; i < shardCount; i++ {
		shard := m.shards[i]

		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}

	return count
}

func (m *ConcurrentMap) Keys() []string {
	count := m.Count()
	ch := make(chan string, count)
	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(shardCount)
		for _, shard := range m.shards {
			go func(shard *ConcurrentMapShard) {
				shard.RLock()
				defer shard.RUnlock()
				defer wg.Done()

				// do
				for key := range shard.items {
					ch <- key
				}
			}(shard)
		}

		wg.Wait()
	}()

	keys := make([]string, 0, count) // len 0
	for key := range ch {
		keys = append(keys, key)
	}

	return keys
}

func (m *ConcurrentMap) getShard(key string) *ConcurrentMapShard {
	return m.shards[uint(fnv32(key))%uint(shardCount)]
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

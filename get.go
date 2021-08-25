package cmap

func (m *ConcurrentMap) Get(key string) (interface{}, bool) {
	shard := m.getShard(key)

	shard.RLock()
	defer shard.RUnlock()
	val, ok := shard.items[key]
	return val, ok
}

func (m *ConcurrentMap) Has(key string) bool {
	shard := m.getShard(key)

	shard.RLock()
	defer shard.RUnlock()
	_, ok := shard.items[key]
	return ok
}

func (m *ConcurrentMap) Pop(key string) (interface{}, bool) {
	shard := m.getShard(key)

	shard.Lock()
	defer shard.Unlock()
	val, ok := shard.items[key]
	if !ok {
		return nil, ok
	}

	delete(shard.items, key)
	return val, true
}

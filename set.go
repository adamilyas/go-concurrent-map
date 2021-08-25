package cmap

func (m *ConcurrentMap) Set(key string, value interface{}) {
	shard := m.getShard(key)

	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

func (m *ConcurrentMap) SetIfAbsent(key string, value interface{}) bool {
	shard := m.getShard(key)

	shard.Lock()
	defer shard.Unlock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}

	return !ok
}

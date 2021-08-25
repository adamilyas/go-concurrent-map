package cmap

import "reflect"

func (m *ConcurrentMap) Remove(key string) {
	shard := m.getShard(key)

	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

func (m *ConcurrentMap) RemoveIfValue(key string, value interface{}) (ret bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
			return
		}
	}()

	shard := m.getShard(key)

	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)

	if foundValue, ok := shard.items[key]; ok {
		if reflect.DeepEqual(value, foundValue) {
			delete(shard.items, key)
			return true, nil
		}
	}

	return false, nil
}

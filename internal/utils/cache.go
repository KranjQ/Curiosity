package utils

import "sync"

type Cache[K comparable, V any] struct {
	items map[K]V
	mu    *sync.RWMutex
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]V),
		mu:    &sync.RWMutex{},
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	c.items[key] = value
	c.mu.Unlock()
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, exists := c.items[key]
	if !exists {
		var zeroValue V
		return zeroValue, false
	}
	return item, true
}

func (c *Cache[K, V]) GetAll() []V {
	var values []V

	for _, val := range c.items {
		values = append(values, val)
	}
	return values
}

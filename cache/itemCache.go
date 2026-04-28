package cache

import (
	"sync"
)

type itemCacheItem[V any] struct {
	Value V
}

type ItemCache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]itemCacheItem[V]
}

func NewItemCache[K comparable, V any]() *ItemCache[K, V] {
	return &ItemCache[K, V]{
		mu:    sync.RWMutex{},
		items: make(map[K]itemCacheItem[V]),
	}
}

func (c *ItemCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = itemCacheItem[V]{
		Value: value,
	}
}

func (c *ItemCache[K, V]) Get(key K) V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return item.Value
	}

	return item.Value
}

func (c *ItemCache[K, V]) GetAll() map[K]V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[K]V)
	for key, item := range c.items {
		result[key] = item.Value
	}

	return result
}

func (c *ItemCache[K, V]) GetOk(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return item.Value, false
	}

	return item.Value, true
}

func (c *ItemCache[K, V]) RemovePosition(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *ItemCache[K, V]) PopAll() map[K]V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[K]V)
	for key, item := range c.items {
		result[key] = item.Value
	}

	c.items = make(map[K]itemCacheItem[V])

	return result
}

func (c *ItemCache[K, V]) ChangeOrCreate(key K, f func(V) V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ищем запись по ключу
	item, found := c.items[key]

	// Если запись не найдена
	if !found {

		// Создаем новую запись
		var emptyType V

		c.items[key] = itemCacheItem[V]{
			Value: f(emptyType),
		}

	} else { // Если найдена

		// Обновляем запись
		c.items[key] = itemCacheItem[V]{
			Value: f(item.Value),
		}
	}
}

package cache

import (
	"sync"
	"time"
)

type itemCacheItem[V any] struct {
	Value      V
	Expiration int64
}

type ItemCache[K comparable, V any] struct {
	mu         sync.RWMutex
	items      map[K]itemCacheItem[V]
	defaultTTL time.Duration
}

func NewItemCache[K comparable, V any](defaultTTL time.Duration) *ItemCache[K, V] {
	return &ItemCache[K, V]{
		mu:         sync.RWMutex{},
		items:      make(map[K]itemCacheItem[V]),
		defaultTTL: defaultTTL,
	}
}

func (c *ItemCache[K, V]) Set(key K, value V, ttl ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if len(ttl) > 0 {
		expiration = time.Now().Add(ttl[0]).UnixNano()
	} else {
		expiration = time.Now().Add(c.defaultTTL).UnixNano()
	}

	c.items[key] = itemCacheItem[V]{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *ItemCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || time.Now().UnixNano() > item.Expiration {
		var zero V
		return zero, false
	}

	return item.Value, true
}

func (c *ItemCache[K, V]) PopAll() map[K]V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now().UnixNano()
	result := make(map[K]V)
	for key, item := range c.items {
		if now < item.Expiration {
			result[key] = item.Value
		}
	}

	c.items = make(map[K]itemCacheItem[V])

	return result
}

func (c *ItemCache[K, V]) ChangeOrCreate(key K, f func(V) V, ttl ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()

	// Ищем запись по ключу
	item, found := c.items[key]

	// Если запись не найдена или просрочена
	if !found || now > item.Expiration {

		// Создаем новую запись
		var expiration int64
		if len(ttl) > 0 {
			expiration = time.Now().Add(ttl[0]).UnixNano()
		} else {
			expiration = time.Now().Add(c.defaultTTL).UnixNano()
		}

		var emptyType V

		c.items[key] = itemCacheItem[V]{
			Value:      f(emptyType),
			Expiration: expiration,
		}

	} else { // Если найдена

		// Обновляем запись
		c.items[key] = itemCacheItem[V]{
			Value:      f(item.Value),
			Expiration: item.Expiration,
		}
	}
}

package cache

import (
	"sync"
	"time"
)

type ListCache[V any] struct {
	mu         sync.RWMutex
	items      []V
	defaultTTL time.Duration
	expiration int64
}

func NewListCache[V any](defaultTTL time.Duration) *ListCache[V] {
	return &ListCache[V]{
		mu:         sync.RWMutex{},
		items:      []V{},
		defaultTTL: defaultTTL,
	}
}

func (c *ListCache[V]) Set(values []V, ttl ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if len(ttl) > 0 {
		expiration = time.Now().Add(ttl[0]).UnixNano()
	} else {
		expiration = time.Now().Add(c.defaultTTL).UnixNano()
	}

	c.items = values
	c.expiration = expiration
}

func (c *ListCache[V]) Get() ([]V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if time.Now().UnixNano() > c.expiration {
		return []V{}, false
	}

	return c.items, true
}

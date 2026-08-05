package cache

import (
	"sync"
	"time"
)

type memoryItem struct {
	value      []byte
	expiration int64
}

// MemoryCache is an in-memory implementation of Cache.
type MemoryCache struct {
	items map[string]memoryItem
	mu    sync.RWMutex
}

// NewMemoryCache creates a new MemoryCache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]memoryItem),
	}
}

// Get retrieves a value by key.
func (c *MemoryCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.expiration {
		// Expired. We could actively delete it here, but let's just return false
		// and it will get overwritten or cleared later.
		return nil, false
	}

	return item.value, true
}

// Set stores a value with a TTL.
func (c *MemoryCache) Set(key string, value []byte, ttl time.Duration) {
	expiration := time.Now().Add(ttl).UnixNano()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = memoryItem{
		value:      value,
		expiration: expiration,
	}
}

// Delete removes a key.
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear empties the cache.
func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]memoryItem)
	return nil
}

package cache

import (
	"time"
)

// Cache defines the interface for a key-value store with TTL.
type Cache interface {
	// Get retrieves a value by key. Returns false if not found or expired.
	Get(key string) ([]byte, bool)
	
	// Set stores a value with the given TTL.
	Set(key string, value []byte, ttl time.Duration)
	
	// Delete removes a key from the cache.
	Delete(key string)
	
	// Clear completely empties the cache.
	Clear() error
}

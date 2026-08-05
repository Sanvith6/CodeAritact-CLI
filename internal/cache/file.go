package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// fileItem is the JSON structure saved to disk.
type fileItem struct {
	Value      []byte `json:"value"`
	Expiration int64  `json:"expiration"`
}

// FileCache implements Cache using the filesystem.
type FileCache struct {
	dir string
	mu  sync.Mutex
}

// NewFileCache creates a new FileCache backed by the specified directory.
func NewFileCache(dir string) (*FileCache, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}
	return &FileCache{
		dir: dir,
	}, nil
}

func (c *FileCache) keyToPath(key string) string {
	// Hash the key to prevent path traversal and weird character issues
	hash := sha256.Sum256([]byte(key))
	filename := hex.EncodeToString(hash[:]) + ".json"
	return filepath.Join(c.dir, filename)
}

// Get retrieves a value by key.
func (c *FileCache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	path := c.keyToPath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var item fileItem
	if err := json.Unmarshal(data, &item); err != nil {
		// Corrupted cache file, delete it
		_ = os.Remove(path)
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		_ = os.Remove(path)
		return nil, false
	}

	return item.Value, true
}

// Set stores a value with a TTL.
func (c *FileCache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(ttl).UnixNano()
	item := fileItem{
		Value:      value,
		Expiration: expiration,
	}

	data, err := json.Marshal(item)
	if err != nil {
		return // Ignore marshal errors
	}

	path := c.keyToPath(key)
	_ = os.WriteFile(path, data, 0600)
}

// Delete removes a key.
func (c *FileCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_ = os.Remove(c.keyToPath(key))
}

// Clear empties the cache directory.
func (c *FileCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			_ = os.Remove(filepath.Join(c.dir, entry.Name()))
		}
	}
	return nil
}

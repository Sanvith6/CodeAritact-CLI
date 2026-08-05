package cache

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runCacheTests(t *testing.T, cache Cache) {
	// Test Set and Get
	cache.Set("key1", []byte("value1"), 1*time.Minute)
	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, []byte("value1"), val)

	// Test Miss
	_, ok = cache.Get("missing")
	assert.False(t, ok)

	// Test Delete
	cache.Delete("key1")
	_, ok = cache.Get("key1")
	assert.False(t, ok)

	// Test Expiration
	cache.Set("key2", []byte("value2"), 10*time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	_, ok = cache.Get("key2")
	assert.False(t, ok) // Should be expired

	// Test Clear
	cache.Set("key3", []byte("value3"), 1*time.Minute)
	err := cache.Clear()
	require.NoError(t, err)
	_, ok = cache.Get("key3")
	assert.False(t, ok)
}

func TestMemoryCache(t *testing.T) {
	c := NewMemoryCache()
	runCacheTests(t, c)
}

func TestFileCache(t *testing.T) {
	tmpDir := t.TempDir()
	c, err := NewFileCache(filepath.Join(tmpDir, "cache"))
	require.NoError(t, err)
	
	runCacheTests(t, c)
}

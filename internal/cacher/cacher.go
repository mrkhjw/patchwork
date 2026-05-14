// Package cacher provides a simple TTL-based in-memory cache for patch
// evaluation results, reducing redundant disk and git operations.
package cacher

import (
	"sync"
	"time"
)

// Entry holds a cached value and its expiry time.
type Entry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// Cache is a thread-safe TTL cache.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]Entry
	default TTL time.Duration
}

// New creates a Cache with the given default TTL.
func New(ttl time.Duration) *Cache {
	return &Cache{
		items:      make(map[string]Entry),
		defaultTTL: ttl,
	}
}

// Set stores a value under key with the default TTL.
func (c *Cache) Set(key string, value interface{}) {
	c.SetTTL(key, value, c.defaultTTL)
}

// SetTTL stores a value under key with an explicit TTL.
func (c *Cache) SetTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a value by key. Returns (value, true) if present and not
// expired, otherwise (nil, false).
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[key]
	if !ok || time.Now().After(e.ExpiresAt) {
		return nil, false
	}
	return e.Value, true
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]Entry)
}

// Evict removes all expired entries and returns the number evicted.
func (c *Cache) Evict() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	count := 0
	for k, e := range c.items {
		if now.After(e.ExpiresAt) {
			delete(c.items, k)
			count++
		}
	}
	return count
}

// Len returns the number of entries currently in the cache (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

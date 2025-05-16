package cache

import (
	"container/list"
	"sync"
	"time"
)

// LRUCache implements a thread-safe LRU cache with time-based expiration
type LRUCache struct {
	capacity    int
	items       map[string]*list.Element
	evictionQ   *list.List
	mu          sync.RWMutex
	ttl         time.Duration
	cleanupTick time.Duration
	stopChan    chan struct{}
}

// cacheItem represents an item in the cache
type cacheItem struct {
	key       string
	value     interface{}
	expiresAt time.Time
}

// NewLRUCache creates a new LRU cache with the given capacity and TTL
func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	cache := &LRUCache{
		capacity:    capacity,
		items:       make(map[string]*list.Element, capacity),
		evictionQ:   list.New(),
		ttl:         ttl,
		cleanupTick: time.Minute,
		stopChan:    make(chan struct{}),
	}

	// Start cleanup routine
	go cache.cleanupRoutine()

	return cache
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	element, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		return nil, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Re-check after acquiring the write lock
	element, exists = c.items[key]
	if !exists {
		return nil, false
	}

	item := element.Value.(*cacheItem)

	// Check if expired
	if time.Now().After(item.expiresAt) {
		c.evictionQ.Remove(element)
		delete(c.items, key)
		return nil, false
	}

	// Move to front as it was recently accessed
	c.evictionQ.MoveToFront(element)
	return item.value, true
}

// Set adds or updates a value in the cache
func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If key exists, update it and move to front
	if element, exists := c.items[key]; exists {
		c.evictionQ.MoveToFront(element)
		item := element.Value.(*cacheItem)
		item.value = value
		item.expiresAt = time.Now().Add(c.ttl)
		return
	}

	// If we're at capacity, remove the least recently used item
	if c.evictionQ.Len() >= c.capacity {
		oldest := c.evictionQ.Back()
		if oldest != nil {
			c.evictionQ.Remove(oldest)
			item := oldest.Value.(*cacheItem)
			delete(c.items, item.key)
		}
	}

	// Add new item
	item := &cacheItem{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
	element := c.evictionQ.PushFront(item)
	c.items[key] = element
}

// Delete removes an item from the cache
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.items[key]; exists {
		c.evictionQ.Remove(element)
		delete(c.items, key)
	}
}

// cleanupRoutine periodically removes expired items
func (c *LRUCache) cleanupRoutine() {
	ticker := time.NewTicker(c.cleanupTick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stopChan:
			return
		}
	}
}

// removeExpired removes all expired items from the cache
func (c *LRUCache) removeExpired() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, element := range c.items {
		item := element.Value.(*cacheItem)
		if now.After(item.expiresAt) {
			c.evictionQ.Remove(element)
			delete(c.items, key)
		}
	}
}

// Close stops the cleanup routine
func (c *LRUCache) Close() {
	close(c.stopChan)
}

// Size returns the current size of the cache
func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

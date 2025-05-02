package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      *sync.RWMutex
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for true {
		c.mu.RLock()
		for key, entry := range c.entries {
			present := time.Now()
			if present.Sub(entry.createdAt) > interval {
				c.mu.RUnlock()
				c.mu.Lock()
				delete(c.entries, key)
				c.mu.Unlock()
				c.mu.RLock()
			}
		}
		c.mu.RUnlock()
		<-ticker.C
	}
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		entries: map[string]cacheEntry{},
		mu:      &sync.RWMutex{},
	}
	go c.reapLoop(interval)
	return c
}

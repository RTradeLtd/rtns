package rtns

import "sync"

// Cache is used to cache keyIDs
// we have published
type Cache struct {
	entries map[string]bool
	mux     sync.RWMutex
}

// NewCache is used to instantiate a new cache
func NewCache() Cache {
	return Cache{make(map[string]bool), sync.RWMutex{}}
}

// Set is used to set an item in cache
func (c Cache) Set(entry string) {
	if c.exists(entry) {
		return
	}
	c.mux.Lock()
	c.entries[entry] = true
	c.mux.Unlock()
}

func (c Cache) exists(entry string) bool {
	c.mux.RLock()
	_, exists := c.entries[entry]
	c.mux.RUnlock()
	return exists
}

// List is used to list items in our cache
func (c Cache) List() []string {
	c.mux.RLock()
	var entries []string
	for entry := range c.entries {
		entries = append(entries, entry)
	}
	c.mux.RUnlock()
	return entries
}

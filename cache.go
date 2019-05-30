package rtns

import "sync"

// cache is used to cache keyIDs
// we have published
type cache struct {
	entries map[string]bool
	mux     sync.RWMutex
}

// newCache is used to instantiate a new cache
func newCache() *cache {
	return &cache{entries: make(map[string]bool)}
}

// set is used to set an item in cache
func (c *cache) Set(entry string) {
	if c.exists(entry) {
		return
	}
	c.mux.Lock()
	c.entries[entry] = true
	c.mux.Unlock()
}

func (c *cache) exists(entry string) bool {
	c.mux.RLock()
	_, exists := c.entries[entry]
	c.mux.RUnlock()
	return exists
}

// list is used to list items in our cache
func (c *cache) list() []string {
	c.mux.RLock()
	var entries []string
	for entry := range c.entries {
		entries = append(entries, entry)
	}
	c.mux.RUnlock()
	return entries
}

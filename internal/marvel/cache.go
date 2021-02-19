package marvel

import "sync"

// Cache is the interface for caching Marvel data.
type Cache interface {
	SetCharacterIds(charIds IntSet)
	GetCharacterIds() IntSet
}

// InMemCache is the in-memory implementation of Cache.
type InMemCache struct {
	sync.RWMutex
	characterIds *IntSet
}

func NewInMemCache() *InMemCache {
	return &InMemCache{
		characterIds: NewIntSet(),
	}
}

// SetCharacterIds replaces the cached set of character IDs with `charIds`.
func (c *InMemCache) SetCharacterIds(charIds IntSet) {
	c.Lock()
	defer c.Unlock()
	c.characterIds = &charIds
}

// GetCharacterIds returns the cached set of character IDs.
func (c *InMemCache) GetCharacterIds() IntSet {
	c.RLock()
	defer c.RUnlock()
	return *c.characterIds
}

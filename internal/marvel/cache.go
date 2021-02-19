package marvel

import "sync"

type Cache interface {
	SetCharacterIds(charIds IntSet)
	GetCharacterIds() IntSet
}

type InMemCache struct {
	sync.RWMutex
	characterIds *IntSet
}

func NewInMemCache() *InMemCache {
	return &InMemCache{
		characterIds: NewIntSet(),
	}
}

func (c *InMemCache) SetCharacterIds(charIds IntSet) {
	c.Lock()
	defer c.Unlock()
	c.characterIds = &charIds
}

func (c *InMemCache) GetCharacterIds() IntSet {
	c.RLock()
	defer c.RUnlock()
	return *c.characterIds
}

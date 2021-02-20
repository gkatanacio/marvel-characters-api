package marvel

import (
	"sync"
	"time"
)

// CharacterCache is the interface for caching Marvel character data.
type CharacterCache interface {
	SetCharacterIds(charIds IntSet, latestModified time.Time)
	GetCharacterIds() (IntSet, time.Time)
}

// InMemCharacterCache is the in-memory implementation of CharacterCache.
type InMemCharacterCache struct {
	sync.RWMutex
	characterIds   *IntSet
	latestModified time.Time
}

func NewInMemCache() *InMemCharacterCache {
	return &InMemCharacterCache{
		characterIds: NewIntSet(),
	}
}

// SetCharacterIds replaces the cached set of character IDs with `charIds` and
// the corresponding latest modified time with `latestModified`.
func (c *InMemCharacterCache) SetCharacterIds(charIds IntSet, latestModified time.Time) {
	c.Lock()
	defer c.Unlock()
	c.characterIds = &charIds
	c.latestModified = latestModified
}

// GetCharacterIds returns the cached set of character IDs along with the
// corresponding latest modified time.
func (c *InMemCharacterCache) GetCharacterIds() (IntSet, time.Time) {
	c.RLock()
	defer c.RUnlock()
	return *c.characterIds, c.latestModified
}

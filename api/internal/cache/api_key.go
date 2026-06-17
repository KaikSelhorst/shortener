package cache

import (
	"sync"
	"time"

	"github.com/KaikSelhorst/shortener/internal/model"
)

type apiKeyCacheEntry struct {
	key       *model.APIKey
	expiresAt time.Time
}

// APIKeyCache is a short-lived in-memory cache for API key lookups.
// Reduces database round-trips on high-frequency authenticated requests.
// TTL is intentionally short so revoked keys stop working quickly.
type APIKeyCache struct {
	mu    sync.RWMutex
	items map[string]*apiKeyCacheEntry
	ttl   time.Duration
}

func NewAPIKeyCache(ttl time.Duration) *APIKeyCache {
	return &APIKeyCache{
		items: make(map[string]*apiKeyCacheEntry),
		ttl:   ttl,
	}
}

func (c *APIKeyCache) Get(hash string) (*model.APIKey, bool) {
	c.mu.RLock()
	e, ok := c.items[hash]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.key, true
}

func (c *APIKeyCache) Set(hash string, key *model.APIKey) {
	c.mu.Lock()
	c.items[hash] = &apiKeyCacheEntry{key: key, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

// Delete evicts a hash immediately — call when an API key is revoked.
func (c *APIKeyCache) Delete(hash string) {
	c.mu.Lock()
	delete(c.items, hash)
	c.mu.Unlock()
}

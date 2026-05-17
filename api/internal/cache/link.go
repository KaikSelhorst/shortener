package cache

import (
	"container/list"
	"sync"
	"time"

	"github.com/KaikSelhorst/shortener/internal/model"
)

type entry struct {
	link      *model.Link
	expiresAt time.Time
	elem      *list.Element
}

// LinkCache is an LRU cache with TTL for short links.
// On capacity overflow the least-recently-used entry is evicted.
// Expired entries are swept in the background every ttl/2 interval.
type LinkCache struct {
	mu      sync.Mutex
	items   map[string]*entry
	lru     *list.List // front = most recently used
	maxSize int
	ttl     time.Duration
	stop    chan struct{}
}

func NewLinkCache(maxSize int, ttl time.Duration) *LinkCache {
	c := &LinkCache{
		items:   make(map[string]*entry, maxSize),
		lru:     list.New(),
		maxSize: maxSize,
		ttl:     ttl,
		stop:    make(chan struct{}),
	}
	go c.cleanup()
	return c
}

func (c *LinkCache) Get(code string) (*model.Link, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.items[code]
	if !ok {
		return nil, false
	}
	if time.Now().After(e.expiresAt) {
		c.remove(code, e)
		return nil, false
	}

	c.lru.MoveToFront(e.elem)
	return e.link, true
}

func (c *LinkCache) Set(link *model.Link) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.items[link.ShortCode]; ok {
		e.link = link
		e.expiresAt = time.Now().Add(c.ttl)
		c.lru.MoveToFront(e.elem)
		return
	}

	if c.lru.Len() >= c.maxSize {
		c.evictLRU()
	}

	elem := c.lru.PushFront(link.ShortCode)
	c.items[link.ShortCode] = &entry{
		link:      link,
		expiresAt: time.Now().Add(c.ttl),
		elem:      elem,
	}
}

func (c *LinkCache) Delete(code string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.items[code]; ok {
		c.remove(code, e)
	}
}

func (c *LinkCache) Close() {
	close(c.stop)
}

func (c *LinkCache) remove(code string, e *entry) {
	c.lru.Remove(e.elem)
	delete(c.items, code)
}

func (c *LinkCache) evictLRU() {
	back := c.lru.Back()
	if back == nil {
		return
	}
	code := back.Value.(string)
	if e, ok := c.items[code]; ok {
		c.remove(code, e)
	}
}

func (c *LinkCache) cleanup() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for code, e := range c.items {
				if now.After(e.expiresAt) {
					c.remove(code, e)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}

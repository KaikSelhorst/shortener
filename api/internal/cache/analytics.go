package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/KaikSelhorst/shortener/internal/model"
)

type analyticsEntry[T any] struct {
	value     T
	expiresAt time.Time
}

// AnalyticsCache caches analytics results by a string key with a fixed TTL.
// Analytics data is expensive to compute but tolerates a few minutes of staleness.
type AnalyticsCache struct {
	mu       sync.RWMutex
	links    map[string]analyticsEntry[*model.LinkAnalytics]
	projects map[string]analyticsEntry[*model.ProjectAnalytics]
	ttl      time.Duration
}

func NewAnalyticsCache(ttl time.Duration) *AnalyticsCache {
	return &AnalyticsCache{
		links:    make(map[string]analyticsEntry[*model.LinkAnalytics]),
		projects: make(map[string]analyticsEntry[*model.ProjectAnalytics]),
		ttl:      ttl,
	}
}

func (c *AnalyticsCache) GetLink(key string) (*model.LinkAnalytics, bool) {
	c.mu.RLock()
	e, ok := c.links[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

func (c *AnalyticsCache) SetLink(key string, v *model.LinkAnalytics) {
	c.mu.Lock()
	c.links[key] = analyticsEntry[*model.LinkAnalytics]{value: v, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *AnalyticsCache) GetProject(key string) (*model.ProjectAnalytics, bool) {
	c.mu.RLock()
	e, ok := c.projects[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

func (c *AnalyticsCache) SetProject(key string, v *model.ProjectAnalytics) {
	c.mu.Lock()
	c.projects[key] = analyticsEntry[*model.ProjectAnalytics]{value: v, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

var analyticsPeriods = []string{"7d", "30d", "90d", ""}

// InvalidateProject removes all period variants for a project from the cache.
func (c *AnalyticsCache) InvalidateProject(projectID int64) {
	c.mu.Lock()
	for _, p := range analyticsPeriods {
		delete(c.projects, fmt.Sprintf("project:%d:%s", projectID, p))
	}
	c.mu.Unlock()
}

// InvalidateLink removes all period variants for a link from the cache.
func (c *AnalyticsCache) InvalidateLink(linkID int64) {
	c.mu.Lock()
	for _, p := range analyticsPeriods {
		delete(c.links, fmt.Sprintf("link:%d:%s", linkID, p))
	}
	c.mu.Unlock()
}

// Close releases all cached entries. Consistent with LinkCache.Close().
func (c *AnalyticsCache) Close() {
	c.mu.Lock()
	c.links = make(map[string]analyticsEntry[*model.LinkAnalytics])
	c.projects = make(map[string]analyticsEntry[*model.ProjectAnalytics])
	c.mu.Unlock()
}

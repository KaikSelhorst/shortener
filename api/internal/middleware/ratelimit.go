package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/KaikSelhorst/shortener/internal/httputil"
)

type windowEntry struct {
	count     int
	windowEnd time.Time
}

// RateLimiter is a fixed-window per-IP rate limiter backed by the standard library.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*windowEntry
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*windowEntry),
		limit:   limit,
		window:  window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for range time.Tick(rl.window) {
		rl.mu.Lock()
		now := time.Now()
		for k, e := range rl.entries {
			if now.After(e.windowEnd) {
				delete(rl.entries, k)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	e, ok := rl.entries[key]
	if !ok || now.After(e.windowEnd) {
		rl.entries[key] = &windowEntry{count: 1, windowEnd: now.Add(rl.window)}
		return true
	}
	if e.count >= rl.limit {
		return false
	}
	e.count++
	return true
}

// Limit is an http.Handler middleware that enforces the rate limit per client IP.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(clientIP(r)) {
			httputil.WriteError(w, http.StatusTooManyRequests, "too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	if i := strings.LastIndex(r.RemoteAddr, ":"); i != -1 {
		return r.RemoteAddr[:i]
	}
	return r.RemoteAddr
}

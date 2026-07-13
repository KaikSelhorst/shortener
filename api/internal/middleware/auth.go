package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/KaikSelhorst/shortener/internal/cache"
	"github.com/KaikSelhorst/shortener/internal/httputil"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	APIKeyKey    contextKey = "api_key"
	requestIDKey contextKey = "request_id"
)

func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(UserIDKey).(int64)
	return id, ok
}

func APIKeyFromContext(ctx context.Context) (*model.APIKey, bool) {
	k, ok := ctx.Value(APIKeyKey).(*model.APIKey)
	return k, ok
}

// ProjectAllowed returns false when the request is authenticated via an API Key
// that is restricted to a specific project and the given projectID doesn't match.
func ProjectAllowed(ctx context.Context, projectID int64) bool {
	key, isAPIKey := APIKeyFromContext(ctx)
	if !isAPIKey || key.ProjectID == nil {
		return true
	}
	return *key.ProjectID == projectID
}

// lastUsedCache tracks the last time each API key ID was persisted to the DB.
// Updates are debounced to at most once per minute per key.
var (
	lastUsedCache sync.Map     // map[int64]time.Time
	lastUsedSem   = make(chan struct{}, 20)
)

func RequireAuth(authService *service.AuthService, apiKeyRepo repository.APIKeyRepo, keyCache *cache.APIKeyCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			if strings.HasPrefix(token, "sk_") {
				hash := service.HashToken(token)

				key, ok := keyCache.Get(hash)
				if !ok {
					var err error
					key, err = apiKeyRepo.GetByHash(r.Context(), hash)
					if err != nil {
						httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
						return
					}
					keyCache.Set(hash, key)
				}

				go updateLastUsed(apiKeyRepo, key.ID)

				ctx := context.WithValue(r.Context(), UserIDKey, key.UserID)
				ctx = context.WithValue(ctx, APIKeyKey, key)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			userID, err := authService.ValidateAccessToken(token)
			if err != nil {
				httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// updateLastUsed persists last_used_at for an API key at most once per minute,
// using a bounded semaphore to cap concurrent DB writes.
func updateLastUsed(repo repository.APIKeyRepo, id int64) {
	now := time.Now()
	if v, ok := lastUsedCache.Load(id); ok && now.Sub(v.(time.Time)) < time.Minute {
		return
	}
	select {
	case lastUsedSem <- struct{}{}:
	default:
		return
	}
	defer func() { <-lastUsedSem }()
	lastUsedCache.Store(id, now)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := repo.UpdateLastUsed(ctx, id); err != nil {
		log.Printf("middleware: failed to update api key last_used_at (id=%d): %v", id, err)
	}
}

// RequireScope checks whether the request has the required scope.
// JWT-authenticated requests implicitly have all scopes.
func RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key, isAPIKey := APIKeyFromContext(r.Context())
			if !isAPIKey {
				next.ServeHTTP(w, r)
				return
			}
			if !key.HasScope(scope) {
				httputil.WriteError(w, http.StatusForbidden, "forbidden: insufficient scope")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

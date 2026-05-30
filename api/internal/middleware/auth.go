package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KaikSelhorst/shortener/internal/httputil"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
)

type contextKey string

const (
	UserIDKey  contextKey = "user_id"
	APIKeyKey  contextKey = "api_key"
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

func RequireAuth(authService *service.AuthService, apiKeyRepo repository.APIKeyRepo) func(http.Handler) http.Handler {
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
				key, err := apiKeyRepo.GetByHash(r.Context(), hash)
				if err != nil {
					httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
					return
				}
				go func(id int64) {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					if err := apiKeyRepo.UpdateLastUsed(ctx, id); err != nil {
						log.Printf("middleware: failed to update api key last_used_at (id=%d): %v", id, err)
					}
				}(key.ID)
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

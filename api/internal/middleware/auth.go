package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	ApiKeyKey contextKey = "api_key"
)

func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(UserIDKey).(int64)
	return id, ok
}

func ApiKeyFromContext(ctx context.Context) (*model.ApiKey, bool) {
	k, ok := ctx.Value(ApiKeyKey).(*model.ApiKey)
	return k, ok
}

// ProjectAllowed returns false when the request is authenticated via an API Key
// that is restricted to a specific project and the given projectID doesn't match.
func ProjectAllowed(ctx context.Context, projectID int64) bool {
	key, isApiKey := ApiKeyFromContext(ctx)
	if !isApiKey || key.ProjectID == nil {
		return true
	}
	return *key.ProjectID == projectID
}

func RequireAuth(authService *service.AuthService, apiKeyRepo *repository.ApiKeyRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			if strings.HasPrefix(token, "sk_") {
				hash := service.HashToken(token)
				key, err := apiKeyRepo.GetByHash(r.Context(), hash)
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				go apiKeyRepo.UpdateLastUsed(context.Background(), key.ID)
				ctx := context.WithValue(r.Context(), UserIDKey, key.UserID)
				ctx = context.WithValue(ctx, ApiKeyKey, key)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			userID, err := authService.ValidateAccessToken(token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireScope verifica se a requisição tem o escopo necessário.
// Requisições autenticadas via JWT têm todos os escopos implicitamente.
func RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key, isApiKey := ApiKeyFromContext(r.Context())
			if !isApiKey {
				next.ServeHTTP(w, r)
				return
			}
			if !key.HasScope(scope) {
				http.Error(w, "Forbidden: insufficient scope", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

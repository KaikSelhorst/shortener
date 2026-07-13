package router

import (
	"context"
	"net/http"
	"time"

	"github.com/KaikSelhorst/shortener/internal/cache"
	"github.com/KaikSelhorst/shortener/internal/config"
	"github.com/KaikSelhorst/shortener/internal/docs"
	"github.com/KaikSelhorst/shortener/internal/handler"
	authmw "github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"go.uber.org/zap"
)

type Handlers struct {
	HealthHandler    *handler.HealthHandler
	RedirectHandler  *handler.RedirectHandler
	ProjectHandler   *handler.ProjectHandler
	LinkHandler      *handler.LinkHandler
	AuthHandler      *handler.AuthHandler
	APIKeyHandler    *handler.APIKeyHandler
	TOTPHandler      *handler.TOTPHandler
	AnalyticsHandler *handler.AnalyticsHandler
	SSEHandler       *handler.SSEHandler
	WebhookHandler   *handler.WebhookHandler
}

type Router struct {
	Server *http.Server
}

// stack chains middlewares: first is outermost, last is innermost.
func stack(mws ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			h = mws[i](h)
		}
		return h
	}
}

func New(cfg *config.Config, handlers *Handlers, authService *service.AuthService, apiKeyRepo repository.APIKeyRepo, keyCache *cache.APIKeyCache, logger *zap.SugaredLogger) *Router {
	mux := http.NewServeMux()

	requireAuth := authmw.RequireAuth(authService, apiKeyRepo, keyCache)
	authLimiter := authmw.NewRateLimiter(10, time.Minute)
	scope := authmw.RequireScope
	fn := func(f func(http.ResponseWriter, *http.Request)) http.Handler { return http.HandlerFunc(f) }

	base := stack(
		authmw.RequestID,
		authmw.NewLogger(logger),
		authmw.SecurityHeaders,
		authmw.RequireJSON,
		func(next http.Handler) http.Handler { return http.MaxBytesHandler(next, 1<<20) },
	)

	auth := stack(base, requireAuth)
	limited := stack(base, authLimiter.Limit)

	mux.Handle("GET /health", base(fn(handlers.HealthHandler.Ok)))
	mux.Handle("GET /openapi.json", base(fn(docs.Handler)))
	mux.Handle("GET /{code}", base(fn(handlers.RedirectHandler.HandleRedirect)))

	mux.Handle("POST /auth/register", limited(fn(handlers.AuthHandler.Register)))
	mux.Handle("POST /auth/login", limited(fn(handlers.AuthHandler.Login)))
	mux.Handle("POST /auth/refresh", base(fn(handlers.AuthHandler.Refresh)))
	mux.Handle("POST /auth/logout", base(fn(handlers.AuthHandler.Logout)))
	mux.Handle("GET /auth/me", auth(fn(handlers.AuthHandler.Me)))
	mux.Handle("POST /auth/mfa/totp", limited(fn(handlers.TOTPHandler.ValidateMFA)))
	mux.Handle("POST /auth/totp/setup", auth(fn(handlers.TOTPHandler.Setup)))
	mux.Handle("POST /auth/totp/confirm", auth(fn(handlers.TOTPHandler.Confirm)))
	mux.Handle("DELETE /auth/totp", auth(fn(handlers.TOTPHandler.Disable)))

	mux.Handle("GET /api-keys", auth(fn(handlers.APIKeyHandler.ListAPIKeys)))
	mux.Handle("POST /api-keys", auth(fn(handlers.APIKeyHandler.CreateAPIKey)))
	mux.Handle("DELETE /api-keys/{id}", auth(fn(handlers.APIKeyHandler.DeleteAPIKey)))

	mux.Handle("GET /projects", auth(scope("projects:read")(fn(handlers.ProjectHandler.ListProjects))))
	mux.Handle("POST /projects", auth(scope("projects:create")(fn(handlers.ProjectHandler.CreateProject))))
	mux.Handle("PUT /projects/{slug}", auth(scope("projects:update")(fn(handlers.ProjectHandler.UpdateProject))))
	mux.Handle("DELETE /projects/{slug}", auth(scope("projects:delete")(fn(handlers.ProjectHandler.DeleteProject))))

	mux.Handle("POST /projects/{slug}/links", auth(scope("links:create")(fn(handlers.LinkHandler.CreateLink))))
	mux.Handle("GET /projects/{slug}/links", auth(scope("links:read")(fn(handlers.LinkHandler.ListLinks))))
	mux.Handle("GET /projects/{slug}/links/{code}", auth(scope("links:read")(fn(handlers.LinkHandler.GetLink))))
	mux.Handle("PUT /projects/{slug}/links/{code}", auth(scope("links:update")(fn(handlers.LinkHandler.UpdateLink))))
	mux.Handle("DELETE /projects/{slug}/links/{code}", auth(scope("links:delete")(fn(handlers.LinkHandler.DeleteLink))))

	mux.Handle("GET /projects/{slug}/analytics", auth(fn(handlers.AnalyticsHandler.GetProjectAnalytics)))
	mux.Handle("GET /projects/{slug}/links/{code}/analytics", auth(fn(handlers.AnalyticsHandler.GetLinkAnalytics)))

	mux.Handle("GET /projects/stream", auth(fn(handlers.SSEHandler.HandleUserStream)))
	mux.Handle("GET /projects/{slug}/stream", auth(fn(handlers.SSEHandler.HandleProjectStream)))
	mux.Handle("GET /projects/{slug}/links/{code}/stream", auth(fn(handlers.SSEHandler.HandleLinkStream)))

	mux.Handle("GET /projects/{slug}/webhooks", auth(scope("webhooks:read")(fn(handlers.WebhookHandler.ListWebhooks))))
	mux.Handle("POST /projects/{slug}/webhooks", auth(scope("webhooks:create")(fn(handlers.WebhookHandler.CreateWebhook))))
	mux.Handle("DELETE /projects/{slug}/webhooks/{id}", auth(scope("webhooks:delete")(fn(handlers.WebhookHandler.DeleteWebhook))))
	mux.Handle("GET /projects/{slug}/webhooks/{id}/deliveries", auth(scope("webhooks:read")(fn(handlers.WebhookHandler.ListDeliveries))))

	server := &http.Server{
		Handler:           mux,
		Addr:              ":" + cfg.Port,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return &Router{Server: server}
}

func (r *Router) Run() error {
	return r.Server.ListenAndServe()
}

func (r *Router) Shutdown(ctx context.Context) error {
	return r.Server.Shutdown(ctx)
}

package router

import (
	"context"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/config"
	"github.com/KaikSelhorst/shortener/internal/handler"
	authmw "github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	HealthHandler   *handler.HealthHandler
	RedirectHandler *handler.RedirectHandler
	ProjectHandler  *handler.ProjectHandler
	LinkHandler     *handler.LinkHandler
	AuthHandler     *handler.AuthHandler
	APIKeyHandler   *handler.APIKeyHandler
}

type Router struct {
	Server *http.Server
}

func New(cfg *config.Config, handlers *Handlers, authService *service.AuthService, apiKeyRepo *repository.APIKeyRepository) *Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	requireAuth := authmw.RequireAuth(authService, apiKeyRepo)

	r.Get("/health", handlers.HealthHandler.Ok)
	r.Get("/{code}", handlers.RedirectHandler.HandleRedirect)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handlers.AuthHandler.Register)
		r.Post("/login", handlers.AuthHandler.Login)
		r.Post("/refresh", handlers.AuthHandler.Refresh)
		r.Post("/logout", handlers.AuthHandler.Logout)
	})

	r.Route("/api-keys", func(r chi.Router) {
		r.Use(requireAuth)
		r.Get("/", handlers.APIKeyHandler.ListAPIKeys)
		r.Post("/", handlers.APIKeyHandler.CreateAPIKey)
		r.Delete("/{id}", handlers.APIKeyHandler.DeleteAPIKey)
	})

	r.Route("/projects", func(r chi.Router) {
		r.Use(requireAuth)

		r.With(authmw.RequireScope("projects:read")).Get("/", handlers.ProjectHandler.ListProjects)
		r.With(authmw.RequireScope("projects:create")).Post("/", handlers.ProjectHandler.CreateProject)
		r.With(authmw.RequireScope("projects:update")).Put("/{slug}", handlers.ProjectHandler.UpdateProject)
		r.With(authmw.RequireScope("projects:delete")).Delete("/{slug}", handlers.ProjectHandler.DeleteProject)

		r.With(authmw.RequireScope("links:create")).Post("/{slug}/links", handlers.LinkHandler.CreateLink)
		r.With(authmw.RequireScope("links:read")).Get("/{slug}/links", handlers.LinkHandler.ListLinks)
		r.With(authmw.RequireScope("links:read")).Get("/{slug}/links/{code}", handlers.LinkHandler.GetLink)
		r.With(authmw.RequireScope("links:update")).Put("/{slug}/links/{code}", handlers.LinkHandler.UpdateLink)
		r.With(authmw.RequireScope("links:delete")).Delete("/{slug}/links/{code}", handlers.LinkHandler.DeleteLink)
	})

	server := &http.Server{Handler: r, Addr: ":" + cfg.Port}

	return &Router{Server: server}
}

func (r *Router) Run() error {
	return r.Server.ListenAndServe()
}

func (r *Router) Shutdown(ctx context.Context) error {
	return r.Server.Shutdown(ctx)
}

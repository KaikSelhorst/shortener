package router

import (
	"context"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/config"
	"github.com/KaikSelhorst/shortener/internal/handler"
	authmw "github.com/KaikSelhorst/shortener/internal/middleware"
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
}

type Router struct {
	Server *http.Server
}

func New(cfg *config.Config, handlers *Handlers, authService *service.AuthService) *Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Get("/health", handlers.HealthHandler.Ok)
	r.Get("/{code}", handlers.RedirectHandler.HandleRedirect)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handlers.AuthHandler.Register)
		r.Post("/login", handlers.AuthHandler.Login)
		r.Post("/refresh", handlers.AuthHandler.Refresh)
		r.Post("/logout", handlers.AuthHandler.Logout)
	})

	r.Route("/projects", func(r chi.Router) {
		r.Use(authmw.RequireAuth(authService))

		r.Post("/", handlers.ProjectHandler.CreateProject)
		r.Put("/{slug}", handlers.ProjectHandler.UpdateProject)
		r.Delete("/{slug}", handlers.ProjectHandler.DeleteProject)

		r.Post("/{slug}/links", handlers.LinkHandler.CreateLink)
		r.Get("/{slug}/links", handlers.LinkHandler.ListLinks)
		r.Get("/{slug}/links/{code}", handlers.LinkHandler.GetLink)
		r.Put("/{slug}/links/{code}", handlers.LinkHandler.UpdateLink)
		r.Delete("/{slug}/links/{code}", handlers.LinkHandler.DeleteLink)
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

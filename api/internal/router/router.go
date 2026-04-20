package router

import (
	"context"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/config"
	"github.com/KaikSelhorst/shortener/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	HealthHandler   *handler.HealthHandler
	RedirectHandler *handler.RedirectHandler
	ProjectHandler  *handler.ProjectHandler
	LinkHandler     *handler.LinkHandler
}

type Router struct {
	Server *http.Server
}

func New(cfg *config.Config, handlers *Handlers) *Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Get("/health", handlers.HealthHandler.Ok)

	r.Post("/projects", handlers.ProjectHandler.CreateProject)
	r.Put("/projects/{slug}", handlers.ProjectHandler.UpdateProject)
	r.Delete("/projects/{slug}", handlers.ProjectHandler.DeleteProject)

	r.Post("/projects/{slug}/links", handlers.LinkHandler.CreateLink)
	r.Get("/projects/{slug}/links", handlers.LinkHandler.ListLinks)
	r.Get("/projects/{slug}/links/{code}", handlers.LinkHandler.GetLink)
	r.Put("/projects/{slug}/links/{code}", handlers.LinkHandler.UpdateLink)
	r.Delete("/projects/{slug}/links/{code}", handlers.LinkHandler.DeleteLink)

	r.Get("/{code}", handlers.RedirectHandler.HandleRedirect)

	server := &http.Server{Handler: r, Addr: ":" + cfg.Port}

	return &Router{Server: server}
}

func (r *Router) Run() error {
	return r.Server.ListenAndServe()
}

func (r *Router) Shutdown(ctx context.Context) error {
	return r.Server.Shutdown(ctx)
}

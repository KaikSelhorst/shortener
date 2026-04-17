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
}

type Router struct {
	Server *http.Server
}

func New(cfg *config.Config, handlers *Handlers) *Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Get("/health", handlers.HealthHandler.Ok)
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

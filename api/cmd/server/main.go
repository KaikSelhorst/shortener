package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KaikSelhorst/shortener/internal/config"
	"github.com/KaikSelhorst/shortener/internal/database"
	"github.com/KaikSelhorst/shortener/internal/handler"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/router"
	"github.com/KaikSelhorst/shortener/migrations"
)

func main() {
	logger, err := config.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()

	if err != nil {
		logger.Fatal(err)
	}

	if err := migrations.Run(cfg.DatabaseURL); err != nil {
		logger.Fatal(err)
	}

	db, err := database.NewDatabase(ctx, cfg.DatabaseURL)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	linkRepository := repository.NewLinkRepository(db.Pool)
	projectRepository := repository.NewProjectRepository(db.Pool)

	healthHandler := handler.NewHealthHandler()
	redirectHandler := handler.NewRedirectHandler(linkRepository)
	projectHandler := handler.NewProjectHandler(projectRepository)

	handlers := &router.Handlers{
		HealthHandler:   healthHandler,
		RedirectHandler: redirectHandler,
		ProjectHandler:  projectHandler,
	}

	r := router.New(cfg, handlers)

	go func() {
		if err := r.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := r.Shutdown(shutdownCtx); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Server stopped")
}

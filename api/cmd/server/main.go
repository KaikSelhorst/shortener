package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KaikSelhorst/shortener/internal/cache"
	"github.com/KaikSelhorst/shortener/internal/config"
	"github.com/KaikSelhorst/shortener/internal/database"
	"github.com/KaikSelhorst/shortener/internal/handler"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/router"
	"github.com/KaikSelhorst/shortener/internal/service"
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
	clickRepository := repository.NewClickRepository(db.Pool)
	userRepository := repository.NewUserRepository(db.Pool)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db.Pool)
	apiKeyRepository := repository.NewAPIKeyRepository(db.Pool)

	tracker := service.NewTrackerService(clickRepository, logger)
	authService := service.NewAuthService(cfg.JWTSecret)

	shortcodeService, err := service.NewShortcodeService()
	if err != nil {
		logger.Fatal(err)
	}

	linkCache := cache.NewLinkCache(1000, 5*time.Minute)
	defer linkCache.Close()

	healthHandler := handler.NewHealthHandler()
	redirectHandler := handler.NewRedirectHandler(linkRepository, tracker, linkCache)
	projectHandler := handler.NewProjectHandler(projectRepository)
	linkHandler := handler.NewLinkHandler(linkRepository, projectRepository, shortcodeService, linkCache, cfg.BaseURL, cfg.CursorSecret)
	authHandler := handler.NewAuthHandler(userRepository, refreshTokenRepository, authService)
	apiKeyHandler := handler.NewAPIKeyHandler(apiKeyRepository, projectRepository, authService)
	totpHandler := handler.NewTOTPHandler(userRepository, refreshTokenRepository, authService)
	analyticsCache := cache.NewAnalyticsCache(5 * time.Minute)
	analyticsHandler := handler.NewAnalyticsHandler(clickRepository, projectRepository, linkRepository, analyticsCache)

	handlers := &router.Handlers{
		HealthHandler:    healthHandler,
		RedirectHandler:  redirectHandler,
		ProjectHandler:   projectHandler,
		LinkHandler:      linkHandler,
		AuthHandler:      authHandler,
		APIKeyHandler:    apiKeyHandler,
		TOTPHandler:      totpHandler,
		AnalyticsHandler: analyticsHandler,
	}

	r := router.New(cfg, handlers, authService, apiKeyRepository)

	go func() {
		if err := r.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	tracker.Shutdown()

	logger.Info("Server stopped")
}

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
	"github.com/KaikSelhorst/shortener/internal/router"
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

	cfg, err := config.Load()

	if err != nil {
		logger.Fatal(err)
	}

	r := router.New(cfg)

	go func() {
		if err := r.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := r.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Server stopped")
}

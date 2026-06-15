package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port              string
	DatabaseURL       string
	BaseURL           string
	CursorSecret      string
	JWTSecret         string
	IPHashSecret      string
	WebhookSecretKey  string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	databaseURL := os.Getenv("DATABASE_URL")
	baseURL := os.Getenv("BASE_URL")
	cursorSecret := os.Getenv("CURSOR_SECRET")
	jwtSecret := os.Getenv("JWT_SECRET")
	ipHashSecret := os.Getenv("IP_HASH_SECRET")
	webhookSecretKey := os.Getenv("WEBHOOK_SECRET_KEY")

	if port == "" {
		port = "8080"
	}

	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if baseURL == "" {
		return nil, fmt.Errorf("BASE_URL is required")
	}

	if len(cursorSecret) < 16 {
		return nil, fmt.Errorf("CURSOR_SECRET must be at least 16 characters")
	}

	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	if len(ipHashSecret) < 32 {
		return nil, fmt.Errorf("IP_HASH_SECRET must be at least 32 characters")
	}

	if len(webhookSecretKey) < 32 {
		return nil, fmt.Errorf("WEBHOOK_SECRET_KEY must be at least 32 characters")
	}

	return &Config{
		Port:         port,
		DatabaseURL:  databaseURL,
		BaseURL:      baseURL,
		CursorSecret: cursorSecret,
		JWTSecret:        jwtSecret,
		IPHashSecret:     ipHashSecret,
		WebhookSecretKey: webhookSecretKey,
	}, nil
}

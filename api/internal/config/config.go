package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port         string
	DatabaseURL  string
	BaseURL      string
	CursorSecret string
	JWTSecret    string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	databaseURL := os.Getenv("DATABASE_URL")
	baseURL := os.Getenv("BASE_URL")
	cursorSecret := os.Getenv("CURSOR_SECRET")
	jwtSecret := os.Getenv("JWT_SECRET")

	if port == "" {
		port = "8080"
	}

	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if baseURL == "" {
		return nil, fmt.Errorf("BASE_URL is required")
	}

	if cursorSecret == "" {
		return nil, fmt.Errorf("CURSOR_SECRET is required")
	}

	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return &Config{
		Port:         port,
		DatabaseURL:  databaseURL,
		BaseURL:      baseURL,
		CursorSecret: cursorSecret,
		JWTSecret:    jwtSecret,
	}, nil
}

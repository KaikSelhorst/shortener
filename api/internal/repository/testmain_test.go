//go:build integration

package repository_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KaikSelhorst/shortener/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("integration tests require DATABASE_URL to be set")
	}

	if err := migrations.Run(connStr); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	var err error
	testDB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("create pool: %v", err)
	}

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

// truncate clears all tables between tests to ensure isolation.
func truncate(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(context.Background(),
		"TRUNCATE TABLE api_keys, refresh_tokens, links, projects, users RESTART IDENTITY CASCADE",
	)
	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}
}

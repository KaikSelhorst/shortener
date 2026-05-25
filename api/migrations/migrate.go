package migrations

import (
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var FS embed.FS

func Run(databaseURL string) error {
	source, err := iofs.New(FS, ".")
	if err != nil {
		return fmt.Errorf("migrations source: %w", err)
	}

	pgx5URL := "pgx5://" + strings.TrimPrefix(strings.TrimPrefix(databaseURL, "postgres://"), "postgresql://")
	m, err := migrate.NewWithSourceInstance("iofs", source, pgx5URL)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

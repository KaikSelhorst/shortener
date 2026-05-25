package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(ctx context.Context, databaseUrl string) (*Database, error) {
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	// Cap connections to avoid overloading the database.
	// MinConns is intentionally left at 0 (default) so the pool never forces
	// idle connections open — a positive minimum would keep connections alive
	// that the database server may kill due to its own idle timeout.
	config.MaxConns = 25
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Database{Pool: pool}, nil
}

func (c *Database) Close() {
	c.Pool.Close()
}

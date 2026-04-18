package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(ctx context.Context, databaseUrl string) (*Database, error) {

	pool, err := pgxpool.New(ctx, databaseUrl)

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

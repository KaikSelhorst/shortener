package repository

import (
	"context"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClickRepository struct {
	db *pgxpool.Pool
}

func NewClickRepository(db *pgxpool.Pool) *ClickRepository {
	return &ClickRepository{db: db}
}

func (r *ClickRepository) BatchInsert(ctx context.Context, clicks []model.Click) error {
	rows := make([][]any, len(clicks))
	for i, c := range clicks {
		rows[i] = []any{c.LinkID, c.UserAgent, c.IPAddress, c.Referer}
	}
	_, err := r.db.CopyFrom(ctx,
		pgx.Identifier{"clicks"},
		[]string{"link_id", "user_agent", "ip_address", "referer"},
		pgx.CopyFromRows(rows),
	)
	return err
}

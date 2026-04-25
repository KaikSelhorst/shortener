package repository

import (
	"context"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, rt *model.RefreshToken) error {
	return r.db.QueryRow(ctx,
		"INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id, created_at",
		rt.UserID, rt.TokenHash, rt.ExpiresAt,
	).Scan(&rt.ID, &rt.CreatedAt)
}

func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, hash string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.QueryRow(ctx,
		"SELECT id, user_id, token_hash, expires_at, created_at, revoked_at FROM refresh_tokens WHERE token_hash = $1",
		hash,
	).Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt, &rt.RevokedAt)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx,
		"UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1",
		id,
	)
	return err
}

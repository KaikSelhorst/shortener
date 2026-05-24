package repository

import (
	"context"
	"errors"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKeyRepository struct {
	db *pgxpool.Pool
}

func NewAPIKeyRepository(db *pgxpool.Pool) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) Create(ctx context.Context, key *model.APIKey) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO api_keys (user_id, project_id, name, key_prefix, key_hash, scopes)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created_at`,
		key.UserID, key.ProjectID, key.Name, key.KeyPrefix, key.KeyHash, key.Scopes,
	).Scan(&key.ID, &key.CreatedAt)
}

func (r *APIKeyRepository) List(ctx context.Context, userID int64) ([]*model.APIKey, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, project_id, name, key_prefix, scopes, last_used_at, created_at
		 FROM api_keys WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*model.APIKey
	for rows.Next() {
		var k model.APIKey
		if err := rows.Scan(&k.ID, &k.UserID, &k.ProjectID, &k.Name, &k.KeyPrefix, &k.Scopes, &k.LastUsedAt, &k.CreatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, &k)
	}
	return keys, rows.Err()
}

func (r *APIKeyRepository) GetByHash(ctx context.Context, hash string) (*model.APIKey, error) {
	var k model.APIKey
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, project_id, name, key_prefix, key_hash, scopes, last_used_at, created_at
		 FROM api_keys WHERE key_hash = $1`,
		hash,
	).Scan(&k.ID, &k.UserID, &k.ProjectID, &k.Name, &k.KeyPrefix, &k.KeyHash, &k.Scopes, &k.LastUsedAt, &k.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &k, nil
}

func (r *APIKeyRepository) Delete(ctx context.Context, id int64, userID int64) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM api_keys WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `UPDATE api_keys SET last_used_at = NOW() WHERE id = $1`, id)
	return err
}

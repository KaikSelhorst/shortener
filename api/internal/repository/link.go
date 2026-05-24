package repository

import (
	"context"
	"errors"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) GetByCode(ctx context.Context, code string) (*model.Link, error) {
	var link model.Link
	err := r.db.QueryRow(ctx,
		"SELECT id, project_id, short_code, original_url, title, description, og_image, expires_at, created_at FROM links WHERE short_code = $1", code,
	).Scan(&link.ID, &link.ProjectID, &link.ShortCode, &link.OriginalURL, &link.Title, &link.Description, &link.OgImage, &link.ExpiresAt, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) List(ctx context.Context, projectID int64, cursor uint64, direction string, limit int) ([]*model.Link, error) {
	const cols = "SELECT id, project_id, short_code, original_url, title, description, og_image, expires_at, created_at FROM links"

	var (
		rows pgx.Rows
		err  error
	)

	switch {
	case cursor == 0:
		rows, err = r.db.Query(ctx, cols+" WHERE project_id = $1 ORDER BY id DESC LIMIT $2", projectID, limit)
	case direction == "prev":
		rows, err = r.db.Query(ctx, cols+" WHERE project_id = $1 AND id > $2 ORDER BY id DESC LIMIT $3", projectID, cursor, limit)
	default:
		rows, err = r.db.Query(ctx, cols+" WHERE project_id = $1 AND id < $2 ORDER BY id DESC LIMIT $3", projectID, cursor, limit)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var links []*model.Link
	for rows.Next() {
		var link model.Link
		if err := rows.Scan(&link.ID, &link.ProjectID, &link.ShortCode, &link.OriginalURL, &link.Title, &link.Description, &link.OgImage, &link.ExpiresAt, &link.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, rows.Err()
}

func (r *LinkRepository) Create(ctx context.Context, link *model.Link, generateCode func(uint64) (string, error)) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		"INSERT INTO links (project_id, original_url, title, description, og_image, expires_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at",
		link.ProjectID, link.OriginalURL, link.Title, link.Description, link.OgImage, link.ExpiresAt,
	).Scan(&link.ID, &link.CreatedAt)
	if err != nil {
		return err
	}

	shortCode, err := generateCode(uint64(link.ID))
	if err != nil {
		return err
	}
	link.ShortCode = shortCode

	_, err = tx.Exec(ctx, "UPDATE links SET short_code = $1 WHERE id = $2", shortCode, link.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *LinkRepository) Update(ctx context.Context, link *model.Link) error {
	_, err := r.db.Exec(ctx,
		"UPDATE links SET original_url = $1, title = $2, description = $3, og_image = $4, expires_at = $5 WHERE id = $6",
		link.OriginalURL, link.Title, link.Description, link.OgImage, link.ExpiresAt, link.ID,
	)
	return err
}

func (r *LinkRepository) DeleteByCode(ctx context.Context, code string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM links WHERE short_code = $1", code)
	return err
}

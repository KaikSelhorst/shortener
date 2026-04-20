package repository

import (
	"context"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) GetByID(ctx context.Context, id int) (*model.Link, error) {
	var link model.Link

	err := r.db.QueryRow(ctx,
		"SELECT id, project_id, original_url, title, description, og_image, created_at FROM links WHERE id = $1", id,
	).Scan(&link.ID, &link.ProjectID, &link.OriginalURL, &link.Title, &link.Description, &link.OgImage, &link.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *LinkRepository) Create(ctx context.Context, link *model.Link) error {
	err := r.db.QueryRow(ctx,
		"INSERT INTO links (project_id, original_url, title, description, og_image) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at",
		link.ProjectID, link.OriginalURL, link.Title, link.Description, link.OgImage,
	).Scan(&link.ID, &link.CreatedAt)

	return err
}

func (r *LinkRepository) DeleteByID(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM links WHERE id = $1", id)
	return err
}

func (r *LinkRepository) Update(ctx context.Context, link *model.Link) error {
	_, err := r.db.Exec(ctx,
		"UPDATE links SET project_id = $1, original_url = $2, title = $3, description = $4, og_image = $5 WHERE id = $6",
		link.ProjectID, link.OriginalURL, link.Title, link.Description, link.OgImage, link.ID,
	)

	return err
}

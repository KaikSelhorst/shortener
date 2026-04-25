package repository

import (
	"context"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	var p model.Project
	err := r.db.QueryRow(ctx,
		"SELECT id, user_id, name, slug, created_at FROM projects WHERE id = $1", id,
	).Scan(&p.ID, &p.UserID, &p.Name, &p.Slug, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepository) FindBySlug(ctx context.Context, slug string) (*model.Project, error) {
	var p model.Project
	err := r.db.QueryRow(ctx,
		"SELECT id, user_id, name, slug, created_at FROM projects WHERE slug = $1", slug,
	).Scan(&p.ID, &p.UserID, &p.Name, &p.Slug, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepository) Create(ctx context.Context, project *model.Project) error {
	return r.db.QueryRow(ctx,
		"INSERT INTO projects (user_id, name, slug) VALUES ($1, $2, $3) RETURNING id, created_at",
		project.UserID, project.Name, project.Slug,
	).Scan(&project.ID, &project.CreatedAt)
}

func (r *ProjectRepository) Update(ctx context.Context, project *model.Project) error {
	_, err := r.db.Exec(ctx,
		"UPDATE projects SET name = $1, slug = $2 WHERE id = $3",
		project.Name, project.Slug, project.ID,
	)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM projects WHERE id = $1", id)
	return err
}

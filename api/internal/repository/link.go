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

func (r *LinkRepository) GetByCodeWithStats(ctx context.Context, code string) (*model.Link, error) {
	var link model.Link
	err := r.db.QueryRow(ctx,
		`SELECT l.id, l.project_id, l.short_code, l.original_url, l.title, l.description, l.og_image, l.expires_at, l.created_at, COUNT(c.id) AS total_clicks
FROM links l LEFT JOIN clicks c ON c.link_id = l.id
WHERE l.short_code = $1
GROUP BY l.id`, code,
	).Scan(&link.ID, &link.ProjectID, &link.ShortCode, &link.OriginalURL, &link.Title, &link.Description, &link.OgImage, &link.ExpiresAt, &link.CreatedAt, &link.TotalClicks)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return links, nil
	}

	// Batch-fetch click counts for the returned links in a single targeted query.
	// Splitting from the main query avoids a GROUP BY aggregation over the full
	// clicks table on every list request.
	ids := make([]int64, len(links))
	for i, l := range links {
		ids[i] = l.ID
	}

	countRows, err := r.db.Query(ctx,
		"SELECT link_id, COUNT(*) FROM clicks WHERE link_id = ANY($1) GROUP BY link_id",
		ids,
	)
	if err != nil {
		return nil, err
	}
	defer countRows.Close()

	counts := make(map[int64]int64, len(links))
	for countRows.Next() {
		var linkID, n int64
		if err := countRows.Scan(&linkID, &n); err != nil {
			return nil, err
		}
		counts[linkID] = n
	}
	if err := countRows.Err(); err != nil {
		return nil, err
	}

	for _, l := range links {
		l.TotalClicks = counts[l.ID]
	}

	return links, nil
}

func (r *LinkRepository) Create(ctx context.Context, link *model.Link, generateCode func(uint64) (string, error)) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	// Rollback is a no-op after a successful Commit; the error is intentionally
	// ignored because it is always pgx.ErrTxClosed in that case.
	defer func() { _ = tx.Rollback(ctx) }()

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

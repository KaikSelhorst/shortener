package repository

import (
	"context"
	"time"

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
		rows[i] = []any{c.LinkID, c.UserAgent, c.IPAddress, c.Referer, c.DeviceType, c.ReferrerSource}
	}
	_, err := r.db.CopyFrom(ctx,
		pgx.Identifier{"clicks"},
		[]string{"link_id", "user_agent", "ip_address", "referer", "device_type", "referrer_source"},
		pgx.CopyFromRows(rows),
	)
	return err
}

func (r *ClickRepository) GetLinkAnalytics(ctx context.Context, linkID int64, since, until time.Time) (*model.LinkAnalytics, error) {
	result := &model.LinkAnalytics{LinkID: linkID}

	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*), COUNT(DISTINCT ip_address)
		 FROM clicks
		 WHERE link_id = $1 AND created_at >= $2 AND created_at < $3`,
		linkID, since, until,
	).Scan(&result.TotalClicks, &result.UniqueClicks); err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT DATE_TRUNC('day', created_at AT TIME ZONE 'UTC'), COUNT(*)
		 FROM clicks
		 WHERE link_id = $1 AND created_at >= $2 AND created_at < $3
		 GROUP BY 1 ORDER BY 1`,
		linkID, since, until,
	)
	if err != nil {
		return nil, err
	}
	result.OverTime, err = scanClicksOverTime(rows)
	if err != nil {
		return nil, err
	}

	result.Devices, err = r.queryDevices(ctx,
		`SELECT device_type, COUNT(*) FROM clicks WHERE link_id = $1 AND created_at >= $2 AND created_at < $3 GROUP BY device_type`,
		linkID, since, until,
	)
	if err != nil {
		return nil, err
	}

	result.Referrers, err = r.queryReferrers(ctx,
		`SELECT referrer_source, COUNT(*) FROM clicks WHERE link_id = $1 AND created_at >= $2 AND created_at < $3 GROUP BY referrer_source`,
		linkID, since, until,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ClickRepository) GetProjectAnalytics(ctx context.Context, projectID int64, since, until time.Time) (*model.ProjectAnalytics, error) {
	result := &model.ProjectAnalytics{}

	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*), COUNT(DISTINCT c.ip_address)
		 FROM clicks c
		 JOIN links l ON c.link_id = l.id
		 WHERE l.project_id = $1 AND c.created_at >= $2 AND c.created_at < $3`,
		projectID, since, until,
	).Scan(&result.TotalClicks, &result.UniqueClicks); err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT DATE_TRUNC('day', c.created_at AT TIME ZONE 'UTC'), COUNT(*)
		 FROM clicks c
		 JOIN links l ON c.link_id = l.id
		 WHERE l.project_id = $1 AND c.created_at >= $2 AND c.created_at < $3
		 GROUP BY 1 ORDER BY 1`,
		projectID, since, until,
	)
	if err != nil {
		return nil, err
	}
	result.OverTime, err = scanClicksOverTime(rows)
	if err != nil {
		return nil, err
	}

	result.Devices, err = r.queryDevices(ctx,
		`SELECT c.device_type, COUNT(*)
		 FROM clicks c JOIN links l ON c.link_id = l.id
		 WHERE l.project_id = $1 AND c.created_at >= $2 AND c.created_at < $3
		 GROUP BY c.device_type`,
		projectID, since, until,
	)
	if err != nil {
		return nil, err
	}

	result.Referrers, err = r.queryReferrers(ctx,
		`SELECT c.referrer_source, COUNT(*)
		 FROM clicks c JOIN links l ON c.link_id = l.id
		 WHERE l.project_id = $1 AND c.created_at >= $2 AND c.created_at < $3
		 GROUP BY c.referrer_source`,
		projectID, since, until,
	)
	if err != nil {
		return nil, err
	}

	topRows, err := r.db.Query(ctx,
		`SELECT l.short_code, l.original_url, l.title, COUNT(*) AS clicks
		 FROM clicks c
		 JOIN links l ON c.link_id = l.id
		 WHERE l.project_id = $1 AND c.created_at >= $2 AND c.created_at < $3
		 GROUP BY l.id
		 ORDER BY clicks DESC
		 LIMIT 10`,
		projectID, since, until,
	)
	if err != nil {
		return nil, err
	}
	defer topRows.Close()
	for topRows.Next() {
		var tl model.TopLink
		if err := topRows.Scan(&tl.ShortCode, &tl.OriginalURL, &tl.Title, &tl.TotalClicks); err != nil {
			return nil, err
		}
		result.TopLinks = append(result.TopLinks, tl)
	}
	if err := topRows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func scanClicksOverTime(rows pgx.Rows) ([]model.ClicksOverTime, error) {
	defer rows.Close()
	var out []model.ClicksOverTime
	for rows.Next() {
		var e model.ClicksOverTime
		if err := rows.Scan(&e.Date, &e.Count); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *ClickRepository) queryDevices(ctx context.Context, query string, args ...any) (model.DeviceBreakdown, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.DeviceBreakdown{}, err
	}
	defer rows.Close()

	var d model.DeviceBreakdown
	for rows.Next() {
		var name string
		var count int64
		if err := rows.Scan(&name, &count); err != nil {
			return d, err
		}
		switch name {
		case "mobile":
			d.Mobile = count
		case "desktop":
			d.Desktop = count
		case "tablet":
			d.Tablet = count
		case "bot":
			d.Bot = count
		default:
			d.Unknown += count
		}
	}
	return d, rows.Err()
}

func (r *ClickRepository) queryReferrers(ctx context.Context, query string, args ...any) (model.ReferrerBreakdown, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.ReferrerBreakdown{}, err
	}
	defer rows.Close()

	var rb model.ReferrerBreakdown
	for rows.Next() {
		var name string
		var count int64
		if err := rows.Scan(&name, &count); err != nil {
			return rb, err
		}
		switch name {
		case "direct":
			rb.Direct = count
		case "instagram":
			rb.Instagram = count
		case "facebook":
			rb.Facebook = count
		case "twitter":
			rb.Twitter = count
		case "tiktok":
			rb.TikTok = count
		case "linkedin":
			rb.LinkedIn = count
		case "whatsapp":
			rb.WhatsApp = count
		case "youtube":
			rb.YouTube = count
		case "google":
			rb.Google = count
		default:
			rb.Other += count
		}
	}
	return rb, rows.Err()
}

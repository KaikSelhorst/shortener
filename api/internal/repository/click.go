package repository

import (
	"context"
	"time"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
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
		rows[i] = []any{c.LinkID, c.UserAgent, c.IPHash, c.Referer, c.DeviceType, c.ReferrerSource, c.Browser}
	}
	_, err := r.db.CopyFrom(ctx,
		pgx.Identifier{"clicks"},
		[]string{"link_id", "user_agent", "ip_hash", "referer", "device_type", "referrer_source", "browser"},
		pgx.CopyFromRows(rows),
	)
	return err
}

func (r *ClickRepository) GetLinkAnalytics(ctx context.Context, linkID int64, since, until time.Time) (*model.LinkAnalytics, error) {
	result := &model.LinkAnalytics{
		LinkID:   linkID,
		OverTime: make([]model.ClicksOverTime, 0),
		Browsers: model.BrowserBreakdown{},
	}

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return r.db.QueryRow(gctx,
			`SELECT COUNT(*), COUNT(DISTINCT ip_hash)
			 FROM clicks
			 WHERE link_id = $1 AND created_at >= $2 AND created_at < $3`,
			linkID, since, until,
		).Scan(&result.TotalClicks, &result.UniqueClicks)
	})

	g.Go(func() error {
		rows, err := r.db.Query(gctx,
			`SELECT DATE_TRUNC('day', created_at AT TIME ZONE 'UTC'), COUNT(*)
			 FROM clicks
			 WHERE link_id = $1 AND created_at >= $2 AND created_at < $3
			 GROUP BY 1 ORDER BY 1`,
			linkID, since, until,
		)
		if err != nil {
			return err
		}
		result.OverTime, err = scanClicksOverTime(rows)
		return err
	})

	g.Go(func() error {
		var err error
		result.Devices, err = r.queryDevices(gctx,
			`SELECT device_type, COUNT(*)
			 FROM clicks
			 WHERE link_id = $1 AND created_at >= $2 AND created_at < $3
			 GROUP BY device_type`,
			linkID, since, until,
		)
		return err
	})

	g.Go(func() error {
		var err error
		result.Referrers, err = r.queryReferrers(gctx,
			`SELECT referrer_source, COUNT(*)
			 FROM clicks
			 WHERE link_id = $1 AND created_at >= $2 AND created_at < $3
			 GROUP BY referrer_source`,
			linkID, since, until,
		)
		return err
	})

	g.Go(func() error {
		var err error
		result.Browsers, err = r.queryBrowsers(gctx,
			`SELECT browser, COUNT(*)
			 FROM clicks
			 WHERE link_id = $1 AND created_at >= $2 AND created_at < $3
			 GROUP BY browser`,
			linkID, since, until,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ClickRepository) GetProjectAnalytics(ctx context.Context, projectID int64, since, until time.Time) (*model.ProjectAnalytics, error) {
	result := &model.ProjectAnalytics{
		OverTime: make([]model.ClicksOverTime, 0),
		TopLinks: make([]model.TopLink, 0),
		Browsers: model.BrowserBreakdown{},
	}

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return r.db.QueryRow(gctx,
			`SELECT COUNT(*), COUNT(DISTINCT c.ip_hash)
			 FROM clicks c
			 JOIN links l ON c.link_id = l.id
			 WHERE l.project_id = $1 AND c.created_at >= $2 AND c.created_at < $3`,
			projectID, since, until,
		).Scan(&result.TotalClicks, &result.UniqueClicks)
	})

	g.Go(func() error {
		rows, err := r.db.Query(gctx,
			`SELECT DATE_TRUNC('day', c.created_at AT TIME ZONE 'UTC'), COUNT(*)
			 FROM clicks c
			 JOIN links l ON c.link_id = l.id
			 WHERE l.project_id = $1 AND c.created_at >= $2 AND c.created_at < $3
			 GROUP BY 1 ORDER BY 1`,
			projectID, since, until,
		)
		if err != nil {
			return err
		}
		result.OverTime, err = scanClicksOverTime(rows)
		return err
	})

	g.Go(func() error {
		var err error
		result.Devices, err = r.queryDevices(gctx,
			`SELECT c.device_type, COUNT(*)
			 FROM clicks c JOIN links l ON c.link_id = l.id
			 WHERE l.project_id = $1 AND c.created_at >= $2 AND c.created_at < $3
			 GROUP BY c.device_type`,
			projectID, since, until,
		)
		return err
	})

	g.Go(func() error {
		var err error
		result.Referrers, err = r.queryReferrers(gctx,
			`SELECT c.referrer_source, COUNT(*)
			 FROM clicks c JOIN links l ON c.link_id = l.id
			 WHERE l.project_id = $1 AND c.created_at >= $2 AND c.created_at < $3
			 GROUP BY c.referrer_source`,
			projectID, since, until,
		)
		return err
	})

	g.Go(func() error {
		var err error
		result.Browsers, err = r.queryBrowsers(gctx,
			`SELECT c.browser, COUNT(*)
			 FROM clicks c JOIN links l ON c.link_id = l.id
			 WHERE l.project_id = $1 AND c.created_at >= $2 AND c.created_at < $3
			 GROUP BY c.browser`,
			projectID, since, until,
		)
		return err
	})

	g.Go(func() error {
		rows, err := r.db.Query(gctx,
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
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var tl model.TopLink
			if err := rows.Scan(&tl.ShortCode, &tl.OriginalURL, &tl.Title, &tl.TotalClicks); err != nil {
				return err
			}
			result.TopLinks = append(result.TopLinks, tl)
		}
		return rows.Err()
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ClickRepository) queryBrowsers(ctx context.Context, query string, args ...any) (model.BrowserBreakdown, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.BrowserBreakdown{}, err
	}
	defer rows.Close()

	var b model.BrowserBreakdown
	for rows.Next() {
		var name string
		var count int64
		if err := rows.Scan(&name, &count); err != nil {
			return b, err
		}
		switch name {
		case "chrome":
			b.Chrome = count
		case "firefox":
			b.Firefox = count
		case "safari":
			b.Safari = count
		case "edge":
			b.Edge = count
		case "opera":
			b.Opera = count
		case "samsung":
			b.Samsung = count
		case "ie":
			b.IE = count
		default:
			b.Other += count
		}
	}
	return b, rows.Err()
}

func scanClicksOverTime(rows pgx.Rows) ([]model.ClicksOverTime, error) {
	defer rows.Close()
	out := make([]model.ClicksOverTime, 0)
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
		case "discord":
			rb.Discord = count
		default:
			rb.Other += count
		}
	}
	return rb, rows.Err()
}

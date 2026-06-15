package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WebhookRepo interface {
	Create(ctx context.Context, w *model.Webhook) error
	List(ctx context.Context, projectID int64) ([]*model.Webhook, error)
	GetByID(ctx context.Context, id string, projectID int64) (*model.Webhook, error)
	Delete(ctx context.Context, id string, projectID int64) error
	ListByProjectAndEvent(ctx context.Context, projectID int64, event string) ([]*model.Webhook, error)
	CreateDelivery(ctx context.Context, d *model.WebhookDelivery) error
	ListDeliveries(ctx context.Context, webhookID string, limit int, offset int) ([]*model.WebhookDelivery, error)
	ClaimPendingDeliveries(ctx context.Context, limit int) ([]*PendingDelivery, error)
	UpdateDeliveryStatus(ctx context.Context, id string, status string, responseStatus *int, nextRetryAt *time.Time) error
	ResetStuckDeliveries(ctx context.Context) error
}

// PendingDelivery carries everything the worker needs without an extra fetch.
type PendingDelivery struct {
	ID              string
	WebhookID       string
	Event           string
	Payload         json.RawMessage
	Attempts        int
	TargetURL       string
	EncryptedSecret string
}

type WebhookRepository struct {
	db *pgxpool.Pool
}

func NewWebhookRepository(db *pgxpool.Pool) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) Create(ctx context.Context, w *model.Webhook) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	w.ID = id.String()
	return r.db.QueryRow(ctx,
		`INSERT INTO webhooks (id, project_id, url, secret, events, enabled)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING created_at`,
		w.ID, w.ProjectID, w.URL, w.Secret, w.Events, w.Enabled,
	).Scan(&w.CreatedAt)
}

func (r *WebhookRepository) List(ctx context.Context, projectID int64) ([]*model.Webhook, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, project_id, url, events, enabled, created_at
		 FROM webhooks WHERE project_id = $1 ORDER BY created_at DESC`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []*model.Webhook
	for rows.Next() {
		var w model.Webhook
		if err := rows.Scan(&w.ID, &w.ProjectID, &w.URL, &w.Events, &w.Enabled, &w.CreatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, &w)
	}
	return webhooks, rows.Err()
}

func (r *WebhookRepository) Delete(ctx context.Context, id string, projectID int64) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM webhooks WHERE id = $1 AND project_id = $2`,
		id, projectID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *WebhookRepository) GetByID(ctx context.Context, id string, projectID int64) (*model.Webhook, error) {
	var w model.Webhook
	err := r.db.QueryRow(ctx,
		`SELECT id, project_id, url, events, enabled, created_at
		 FROM webhooks WHERE id = $1 AND project_id = $2`,
		id, projectID,
	).Scan(&w.ID, &w.ProjectID, &w.URL, &w.Events, &w.Enabled, &w.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &w, nil
}

func (r *WebhookRepository) ListByProjectAndEvent(ctx context.Context, projectID int64, event string) ([]*model.Webhook, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, project_id, url, secret, events, enabled, created_at
		 FROM webhooks
		 WHERE project_id = $1 AND enabled = TRUE AND $2 = ANY(events)`,
		projectID, event,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []*model.Webhook
	for rows.Next() {
		var w model.Webhook
		if err := rows.Scan(&w.ID, &w.ProjectID, &w.URL, &w.Secret, &w.Events, &w.Enabled, &w.CreatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, &w)
	}
	return webhooks, rows.Err()
}

func (r *WebhookRepository) CreateDelivery(ctx context.Context, d *model.WebhookDelivery) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO webhook_deliveries (webhook_id, event, payload)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, next_retry_at`,
		d.WebhookID, d.Event, d.Payload,
	).Scan(&d.ID, &d.CreatedAt, &d.NextRetryAt)
}

func (r *WebhookRepository) ListDeliveries(ctx context.Context, webhookID string, limit int, offset int) ([]*model.WebhookDelivery, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, webhook_id, event, payload, status, attempts, response_status, next_retry_at, created_at
		 FROM webhook_deliveries
		 WHERE webhook_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		webhookID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []*model.WebhookDelivery
	for rows.Next() {
		var d model.WebhookDelivery
		if err := rows.Scan(
			&d.ID, &d.WebhookID, &d.Event, &d.Payload,
			&d.Status, &d.Attempts, &d.ResponseStatus, &d.NextRetryAt, &d.CreatedAt,
		); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, &d)
	}
	return deliveries, rows.Err()
}

func (r *WebhookRepository) ClaimPendingDeliveries(ctx context.Context, limit int) ([]*PendingDelivery, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// Step 1: select and lock rows with SKIP LOCKED, joining webhooks for url/secret.
	rows, err := tx.Query(ctx, `
		SELECT wd.id, wd.webhook_id, wd.event, wd.payload, wd.attempts, w.url, w.secret
		FROM webhook_deliveries wd
		JOIN webhooks w ON w.id = wd.webhook_id AND w.enabled = TRUE
		WHERE wd.status = 'pending' AND wd.next_retry_at <= NOW()
		ORDER BY wd.next_retry_at
		LIMIT $1
		FOR UPDATE OF wd SKIP LOCKED`,
		limit,
	)
	if err != nil {
		return nil, err
	}

	var pending []*PendingDelivery
	var ids []string
	for rows.Next() {
		var p PendingDelivery
		if err := rows.Scan(&p.ID, &p.WebhookID, &p.Event, &p.Payload, &p.Attempts, &p.TargetURL, &p.EncryptedSecret); err != nil {
			rows.Close()
			return nil, err
		}
		pending = append(pending, &p)
		ids = append(ids, p.ID)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, tx.Commit(ctx)
	}

	// Step 2: mark claimed rows as processing in the same transaction.
	// The rows are already locked by FOR UPDATE above, so no race is possible.
	_, err = tx.Exec(ctx, `
		UPDATE webhook_deliveries
		SET status = 'processing', attempts = attempts + 1
		WHERE id = ANY($1)`,
		ids,
	)
	if err != nil {
		return nil, err
	}

	return pending, tx.Commit(ctx)
}

// ResetStuckDeliveries moves deliveries that were left in "processing" (e.g.
// after a crash) back to "pending" so they will be retried by the worker.
func (r *WebhookRepository) ResetStuckDeliveries(ctx context.Context) error {
	_, err := r.db.Exec(ctx,
		`UPDATE webhook_deliveries SET status = 'pending', next_retry_at = NOW()
		 WHERE status = 'processing'`,
	)
	return err
}

func (r *WebhookRepository) UpdateDeliveryStatus(ctx context.Context, id string, status string, responseStatus *int, nextRetryAt *time.Time) error {
	var err error
	if nextRetryAt != nil {
		_, err = r.db.Exec(ctx,
			`UPDATE webhook_deliveries SET status = $1, response_status = $2, next_retry_at = $3 WHERE id = $4`,
			status, responseStatus, *nextRetryAt, id,
		)
	} else {
		_, err = r.db.Exec(ctx,
			`UPDATE webhook_deliveries SET status = $1, response_status = $2 WHERE id = $3`,
			status, responseStatus, id,
		)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

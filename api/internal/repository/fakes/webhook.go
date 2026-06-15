package fakes

import (
	"context"
	"sync"
	"time"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/google/uuid"
)

type WebhookRepo struct {
	mu          sync.RWMutex
	webhooks    map[string]*model.Webhook
	deliveries  map[string]*model.WebhookDelivery
	ReturnError error
}

func NewWebhookRepo() *WebhookRepo {
	return &WebhookRepo{
		webhooks:   make(map[string]*model.Webhook),
		deliveries: make(map[string]*model.WebhookDelivery),
	}
}

func (r *WebhookRepo) Create(_ context.Context, w *model.Webhook) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	w.ID = id.String()
	w.CreatedAt = time.Now()
	cp := *w
	r.webhooks[cp.ID] = &cp
	return nil
}

func (r *WebhookRepo) List(_ context.Context, projectID int64) ([]*model.Webhook, error) {
	if r.ReturnError != nil {
		return nil, r.ReturnError
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*model.Webhook
	for _, w := range r.webhooks {
		if w.ProjectID == projectID {
			cp := *w
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (r *WebhookRepo) GetByID(_ context.Context, id string, projectID int64) (*model.Webhook, error) {
	if r.ReturnError != nil {
		return nil, r.ReturnError
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.webhooks[id]
	if !ok || w.ProjectID != projectID {
		return nil, repository.ErrNotFound
	}
	cp := *w
	return &cp, nil
}

func (r *WebhookRepo) Delete(_ context.Context, id string, projectID int64) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.webhooks[id]
	if !ok || w.ProjectID != projectID {
		return repository.ErrNotFound
	}
	delete(r.webhooks, id)
	return nil
}

func (r *WebhookRepo) ListByProjectAndEvent(_ context.Context, projectID int64, event string) ([]*model.Webhook, error) {
	if r.ReturnError != nil {
		return nil, r.ReturnError
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*model.Webhook
	for _, w := range r.webhooks {
		if w.ProjectID != projectID || !w.Enabled {
			continue
		}
		for _, e := range w.Events {
			if e == event {
				cp := *w
				out = append(out, &cp)
				break
			}
		}
	}
	return out, nil
}

func (r *WebhookRepo) CreateDelivery(_ context.Context, d *model.WebhookDelivery) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	d.ID = uuid.New().String()
	d.Status = "pending"
	d.CreatedAt = time.Now()
	cp := *d
	r.deliveries[cp.ID] = &cp
	return nil
}

func (r *WebhookRepo) ListDeliveries(_ context.Context, webhookID string, limit, offset int) ([]*model.WebhookDelivery, error) {
	if r.ReturnError != nil {
		return nil, r.ReturnError
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var all []*model.WebhookDelivery
	for _, d := range r.deliveries {
		if d.WebhookID == webhookID {
			cp := *d
			all = append(all, &cp)
		}
	}
	if offset >= len(all) {
		return nil, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}

func (r *WebhookRepo) ClaimPendingDeliveries(_ context.Context, limit int) ([]*repository.PendingDelivery, error) {
	if r.ReturnError != nil {
		return nil, r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []*repository.PendingDelivery
	for _, d := range r.deliveries {
		if len(out) >= limit {
			break
		}
		if d.Status != "pending" || d.NextRetryAt.After(time.Now()) {
			continue
		}
		w := r.webhooks[d.WebhookID]
		if w == nil || !w.Enabled {
			continue
		}
		d.Status = "processing"
		d.Attempts++
		out = append(out, &repository.PendingDelivery{
			ID:              d.ID,
			WebhookID:       d.WebhookID,
			Event:           d.Event,
			Payload:         d.Payload,
			Attempts:        d.Attempts,
			TargetURL:       w.URL,
			EncryptedSecret: w.Secret,
		})
	}
	return out, nil
}

func (r *WebhookRepo) UpdateDeliveryStatus(_ context.Context, id, status string, responseStatus *int, nextRetryAt *time.Time) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.deliveries[id]
	if !ok {
		return repository.ErrNotFound
	}
	d.Status = status
	d.ResponseStatus = responseStatus
	if nextRetryAt != nil {
		d.NextRetryAt = *nextRetryAt
	}
	return nil
}

func (r *WebhookRepo) ResetStuckDeliveries(_ context.Context) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, d := range r.deliveries {
		if d.Status == "processing" {
			d.Status = "pending"
			d.NextRetryAt = time.Now()
		}
	}
	return nil
}

func (r *WebhookRepo) WebhookCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.webhooks)
}

func (r *WebhookRepo) DeliveryCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.deliveries)
}

//go:build integration

package repository_test

import (
	"errors"
	"testing"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
)

func newWebhookRepo() *repository.WebhookRepository {
	return repository.NewWebhookRepository(testDB)
}

func makeWebhook(t *testing.T, projectID int64, events ...string) *model.Webhook {
	t.Helper()
	if len(events) == 0 {
		events = []string{"link.clicked", "link.created"}
	}
	svc := service.NewWebhookService(nil, []byte("test-key"))
	enc, err := svc.EncryptSecret("plain-secret")
	if err != nil {
		t.Fatalf("encrypt secret: %v", err)
	}
	wh := &model.Webhook{
		ProjectID: projectID,
		URL:       "https://example.com/hook-" + t.Name(),
		Secret:    enc,
		Events:    events,
		Enabled:   true,
	}
	if err := newWebhookRepo().Create(t.Context(), wh); err != nil {
		t.Fatalf("makeWebhook: %v", err)
	}
	return wh
}

// --- Create ---

func TestWebhookRepository_Create_AssignsUUIDAndTimestamp(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)

	wh := &model.Webhook{
		ProjectID: p.ID,
		URL:       "https://example.com/hook",
		Secret:    "encrypted",
		Events:    []string{"link.clicked"},
		Enabled:   true,
	}
	if err := newWebhookRepo().Create(t.Context(), wh); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if wh.ID == "" {
		t.Error("expected non-empty UUID ID")
	}
	if wh.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be populated")
	}
}

// --- List ---

func TestWebhookRepository_List_OnlyReturnsProjectWebhooks(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p1 := makeProject(t, u.ID)
	p2 := makeProject(t, u.ID)

	makeWebhook(t, p1.ID)
	makeWebhook(t, p1.ID)
	makeWebhook(t, p2.ID)

	repo := newWebhookRepo()
	got, err := repo.List(t.Context(), p1.ID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 webhooks for project 1, got %d", len(got))
	}
}

// --- GetByID ---

func TestWebhookRepository_GetByID_Success(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)
	wh := makeWebhook(t, p.ID)

	got, err := newWebhookRepo().GetByID(t.Context(), wh.ID, p.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.ID != wh.ID {
		t.Errorf("expected ID %s, got %s", wh.ID, got.ID)
	}
}

func TestWebhookRepository_GetByID_WrongProject(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)
	wh := makeWebhook(t, p.ID)

	_, err := newWebhookRepo().GetByID(t.Context(), wh.ID, 9999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// --- Delete ---

func TestWebhookRepository_Delete_RemovesWebhook(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)
	wh := makeWebhook(t, p.ID)

	repo := newWebhookRepo()
	if err := repo.Delete(t.Context(), wh.ID, p.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if err := repo.Delete(t.Context(), wh.ID, p.ID); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("expected ErrNotFound on second delete, got %v", err)
	}
}

// --- ListByProjectAndEvent ---

func TestWebhookRepository_ListByProjectAndEvent(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)
	makeWebhook(t, p.ID, "link.clicked", "link.created")
	makeWebhook(t, p.ID, "link.deleted")

	repo := newWebhookRepo()

	got, err := repo.ListByProjectAndEvent(t.Context(), p.ID, "link.clicked")
	if err != nil {
		t.Fatalf("ListByProjectAndEvent: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 webhook for link.clicked, got %d", len(got))
	}

	got, err = repo.ListByProjectAndEvent(t.Context(), p.ID, "link.updated")
	if err != nil {
		t.Fatalf("ListByProjectAndEvent: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 webhooks for unsubscribed event, got %d", len(got))
	}
}

// --- Deliveries ---

func TestWebhookRepository_CreateAndListDeliveries(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)
	wh := makeWebhook(t, p.ID)

	repo := newWebhookRepo()
	for range 3 {
		d := &model.WebhookDelivery{
			WebhookID: wh.ID,
			Event:     "link.clicked",
			Payload:   []byte(`{"event":"link.clicked"}`),
		}
		if err := repo.CreateDelivery(t.Context(), d); err != nil {
			t.Fatalf("CreateDelivery: %v", err)
		}
		if d.ID == "" {
			t.Error("expected delivery ID to be set")
		}
	}

	got, err := repo.ListDeliveries(t.Context(), wh.ID, 10, 0)
	if err != nil {
		t.Fatalf("ListDeliveries: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 deliveries, got %d", len(got))
	}
}

func TestWebhookRepository_ListDeliveries_Pagination(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)
	wh := makeWebhook(t, p.ID)

	repo := newWebhookRepo()
	for range 5 {
		_ = repo.CreateDelivery(t.Context(), &model.WebhookDelivery{
			WebhookID: wh.ID,
			Event:     "link.clicked",
			Payload:   []byte(`{}`),
		})
	}

	page1, _ := repo.ListDeliveries(t.Context(), wh.ID, 3, 0)
	if len(page1) != 3 {
		t.Errorf("expected 3 on page 1, got %d", len(page1))
	}
	page2, _ := repo.ListDeliveries(t.Context(), wh.ID, 3, 3)
	if len(page2) != 2 {
		t.Errorf("expected 2 on page 2, got %d", len(page2))
	}
}

// --- ClaimPendingDeliveries ---

func TestWebhookRepository_ClaimPendingDeliveries(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)
	wh := makeWebhook(t, p.ID)

	repo := newWebhookRepo()
	d := &model.WebhookDelivery{
		WebhookID: wh.ID,
		Event:     "link.clicked",
		Payload:   []byte(`{}`),
	}
	_ = repo.CreateDelivery(t.Context(), d)

	claimed, err := repo.ClaimPendingDeliveries(t.Context(), 10)
	if err != nil {
		t.Fatalf("ClaimPendingDeliveries: %v", err)
	}
	if len(claimed) != 1 {
		t.Fatalf("expected 1 claimed delivery, got %d", len(claimed))
	}
	if claimed[0].WebhookID != wh.ID {
		t.Errorf("unexpected WebhookID: %s", claimed[0].WebhookID)
	}

	// Second claim — should find nothing (already processing)
	claimed2, _ := repo.ClaimPendingDeliveries(t.Context(), 10)
	if len(claimed2) != 0 {
		t.Errorf("expected 0 on second claim, got %d", len(claimed2))
	}
}

// --- ResetStuckDeliveries ---

func TestWebhookRepository_ResetStuckDeliveries(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)
	wh := makeWebhook(t, p.ID)

	repo := newWebhookRepo()
	_ = repo.CreateDelivery(t.Context(), &model.WebhookDelivery{
		WebhookID: wh.ID,
		Event:     "link.clicked",
		Payload:   []byte(`{}`),
	})

	// Claim it → now "processing"
	_, _ = repo.ClaimPendingDeliveries(t.Context(), 10)

	// Reset → back to "pending"
	if err := repo.ResetStuckDeliveries(t.Context()); err != nil {
		t.Fatalf("ResetStuckDeliveries: %v", err)
	}

	claimed, _ := repo.ClaimPendingDeliveries(t.Context(), 10)
	if len(claimed) != 1 {
		t.Errorf("expected delivery to be claimable after reset, got %d claimed", len(claimed))
	}
}

// --- UpdateDeliveryStatus ---

func TestWebhookRepository_UpdateDeliveryStatus(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)
	wh := makeWebhook(t, p.ID)

	repo := newWebhookRepo()
	d := &model.WebhookDelivery{
		WebhookID: wh.ID,
		Event:     "link.clicked",
		Payload:   []byte(`{}`),
	}
	_ = repo.CreateDelivery(t.Context(), d)
	_, _ = repo.ClaimPendingDeliveries(t.Context(), 10)

	code := 200
	if err := repo.UpdateDeliveryStatus(t.Context(), d.ID, "delivered", &code, nil); err != nil {
		t.Fatalf("UpdateDeliveryStatus: %v", err)
	}

	got, _ := repo.ListDeliveries(t.Context(), wh.ID, 10, 0)
	if len(got) == 0 || got[0].Status != "delivered" {
		t.Errorf("expected status 'delivered', got %+v", got)
	}
	if got[0].ResponseStatus == nil || *got[0].ResponseStatus != 200 {
		t.Errorf("expected response_status 200, got %v", got[0].ResponseStatus)
	}
}

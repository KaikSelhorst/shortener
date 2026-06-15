package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/handler"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository/fakes"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/KaikSelhorst/shortener/internal/testutil"
)

func newWebhookHandler() (*handler.WebhookHandler, *fakes.WebhookRepo, *fakes.ProjectRepo) {
	webhooks := fakes.NewWebhookRepo()
	projects := fakes.NewProjectRepo()
	svc := service.NewWebhookService(webhooks, []byte("test-webhook-key"))
	return handler.NewWebhookHandler(svc, webhooks, projects), webhooks, projects
}

func seedWebhook(t *testing.T, repo *fakes.WebhookRepo, projectID int64, url string) *model.Webhook {
	t.Helper()
	wh := &model.Webhook{
		ProjectID: projectID,
		URL:       url,
		Secret:    "encrypted-secret",
		Events:    []string{"link.clicked"},
		Enabled:   true,
	}
	if err := repo.Create(t.Context(), wh); err != nil {
		t.Fatalf("seedWebhook: %v", err)
	}
	return wh
}

func seedDelivery(t *testing.T, repo *fakes.WebhookRepo, webhookID string) {
	t.Helper()
	if err := repo.CreateDelivery(context.Background(), &model.WebhookDelivery{
		WebhookID: webhookID,
		Event:     "link.clicked",
		Payload:   []byte(`{}`),
	}); err != nil {
		t.Fatalf("seedDelivery: %v", err)
	}
}

// --- CreateWebhook ---

func TestWebhookHandler_CreateWebhook_Success(t *testing.T) {
	h, _, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/", dto.CreateWebhookRequest{
		URL:    "https://example.com/hook",
		Events: []string{"link.clicked"},
	})
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug})

	h.CreateWebhook(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp dto.CreateWebhookResponse
	if err := testutil.DecodeJSON(w, &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Error("expected non-empty UUID ID")
	}
	if resp.Secret == "" {
		t.Error("expected plaintext secret in creation response")
	}
	if resp.URL != "https://example.com/hook" {
		t.Errorf("unexpected URL: %s", resp.URL)
	}
}

func TestWebhookHandler_CreateWebhook_InvalidURL(t *testing.T) {
	h, _, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/", dto.CreateWebhookRequest{
		URL:    "not-a-url",
		Events: []string{"link.clicked"},
	})
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug})

	h.CreateWebhook(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWebhookHandler_CreateWebhook_NoEvents(t *testing.T) {
	h, _, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/", dto.CreateWebhookRequest{
		URL:    "https://example.com/hook",
		Events: []string{},
	})
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug})

	h.CreateWebhook(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWebhookHandler_CreateWebhook_InvalidEvent(t *testing.T) {
	h, _, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/", dto.CreateWebhookRequest{
		URL:    "https://example.com/hook",
		Events: []string{"not.an.event"},
	})
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug})

	h.CreateWebhook(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWebhookHandler_CreateWebhook_Forbidden(t *testing.T) {
	h, _, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/", dto.CreateWebhookRequest{
		URL:    "https://example.com/hook",
		Events: []string{"link.clicked"},
	})
	r = testutil.WithUserID(r, 999) // wrong user
	r = withChiParams(r, map[string]string{"slug": p.Slug})

	h.CreateWebhook(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

// --- ListWebhooks ---

func TestWebhookHandler_ListWebhooks_Empty(t *testing.T) {
	h, _, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodGet, "/", nil)
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug})

	h.ListWebhooks(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestWebhookHandler_ListWebhooks_ReturnsOwned(t *testing.T) {
	h, webhooks, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")
	seedWebhook(t, webhooks, p.ID, "https://example.com/hook1")
	seedWebhook(t, webhooks, p.ID, "https://example.com/hook2")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodGet, "/", nil)
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug})

	h.ListWebhooks(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp []dto.WebhookResponse
	if err := testutil.DecodeJSON(w, &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("expected 2 webhooks, got %d", len(resp))
	}
}

// --- DeleteWebhook ---

func TestWebhookHandler_DeleteWebhook_Success(t *testing.T) {
	h, webhooks, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")
	wh := seedWebhook(t, webhooks, p.ID, "https://example.com/hook")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodDelete, "/", nil)
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug, "id": wh.ID})

	h.DeleteWebhook(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
	if webhooks.WebhookCount() != 0 {
		t.Error("expected webhook to be removed from store")
	}
}

func TestWebhookHandler_DeleteWebhook_InvalidUUID(t *testing.T) {
	h, _, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodDelete, "/", nil)
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug, "id": "not-a-uuid"})

	h.DeleteWebhook(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWebhookHandler_DeleteWebhook_NotFound(t *testing.T) {
	h, _, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodDelete, "/", nil)
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{
		"slug": p.Slug,
		"id":   "019401a0-a78e-7a00-a000-000000000001",
	})

	h.DeleteWebhook(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestWebhookHandler_DeleteWebhook_Forbidden(t *testing.T) {
	h, webhooks, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")
	wh := seedWebhook(t, webhooks, p.ID, "https://example.com/hook")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodDelete, "/", nil)
	r = testutil.WithUserID(r, 999) // wrong user
	r = withChiParams(r, map[string]string{"slug": p.Slug, "id": wh.ID})

	h.DeleteWebhook(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

// --- ListDeliveries ---

func TestWebhookHandler_ListDeliveries_Success(t *testing.T) {
	h, webhooks, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")
	wh := seedWebhook(t, webhooks, p.ID, "https://example.com/hook")
	seedDelivery(t, webhooks, wh.ID)
	seedDelivery(t, webhooks, wh.ID)

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodGet, "/?page=1&limit=20", nil)
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug, "id": wh.ID})

	h.ListDeliveries(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp dto.ListDeliveriesResponse
	if err := testutil.DecodeJSON(w, &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Errorf("expected 2 deliveries, got %d", len(resp.Data))
	}
	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
}

func TestWebhookHandler_ListDeliveries_HasMore(t *testing.T) {
	h, webhooks, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")
	wh := seedWebhook(t, webhooks, p.ID, "https://example.com/hook")
	for range 21 {
		seedDelivery(t, webhooks, wh.ID)
	}

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodGet, "/?page=1&limit=20", nil)
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug, "id": wh.ID})

	h.ListDeliveries(w, r)

	var resp dto.ListDeliveriesResponse
	_ = testutil.DecodeJSON(w, &resp)
	if !resp.HasMore {
		t.Error("expected has_more=true with 21 deliveries and limit=20")
	}
	if len(resp.Data) != 20 {
		t.Errorf("expected 20 items in page, got %d", len(resp.Data))
	}
}

func TestWebhookHandler_ListDeliveries_InvalidUUID(t *testing.T) {
	h, _, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodGet, "/", nil)
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{"slug": p.Slug, "id": "bad-id"})

	h.ListDeliveries(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWebhookHandler_ListDeliveries_WebhookNotFound(t *testing.T) {
	h, _, projects := newWebhookHandler()
	p := seedProject(t, projects, 1, "Test", "test-slug")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodGet, "/", nil)
	r = testutil.WithUserID(r, p.UserID)
	r = withChiParams(r, map[string]string{
		"slug": p.Slug,
		"id":   "019401a0-a78e-7a00-a000-000000000001",
	})

	h.ListDeliveries(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

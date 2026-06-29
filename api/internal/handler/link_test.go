package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KaikSelhorst/shortener/internal/cache"
	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/handler"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository/fakes"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/KaikSelhorst/shortener/internal/testutil"
)

func newLinkHandler() (*handler.LinkHandler, *fakes.LinkRepo, *fakes.ProjectRepo) {
	links := fakes.NewLinkRepo()
	projects := fakes.NewProjectRepo()
	svc, _ := service.NewShortcodeService()
	lc := cache.NewLinkCache(100, time.Minute)
	webhookSvc := service.NewWebhookService(fakes.NewWebhookRepo(), []byte("test-key"))
	return handler.NewLinkHandler(links, projects, svc, webhookSvc, lc, "http://short.test", "cursor-secret"), links, projects
}

// withPathValues injects URL path parameters into the request (Go 1.22+).
func withPathValues(r *http.Request, params map[string]string) *http.Request {
	for k, v := range params {
		r.SetPathValue(k, v)
	}
	return r
}

// seedProject creates a project owned by userID in the fake repo.
func seedProject(t *testing.T, projects *fakes.ProjectRepo, userID int64, name, slug string) *model.Project {
	t.Helper()
	p := &model.Project{UserID: userID, Name: name, Slug: slug}
	if err := projects.Create(t.Context(), p); err != nil {
		t.Fatalf("seedProject: %v", err)
	}
	return p
}

// seedLink creates a link in the fake repo using the shortcode service.
func seedLink(t *testing.T, links *fakes.LinkRepo, projectID int64, url string) *model.Link {
	t.Helper()
	svc, _ := service.NewShortcodeService()
	link := &model.Link{ProjectID: projectID, OriginalURL: url}
	if err := links.Create(t.Context(), link, svc.GenerateShortCode); err != nil {
		t.Fatalf("seedLink: %v", err)
	}
	return link
}

// --- CreateLink ---

func TestLinkHandler_CreateLink_Success(t *testing.T) {
	t.Parallel()
	h, _, projects := newLinkHandler()
	const userID = int64(1)
	seedProject(t, projects, userID, "My Project", "my-project")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		withPathValues(
			testutil.NewRequest(http.MethodPost, "/projects/my-project/links", dto.LinkRequest{
				URL: "https://example.com",
			}),
			map[string]string{"slug": "my-project"},
		),
		userID,
	)

	h.CreateLink(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var res dto.LinkResponse
	if err := testutil.DecodeJSON(w, &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if res.ShortCode == "" {
		t.Error("expected short code in response")
	}
	if res.OriginalURL != "https://example.com" {
		t.Errorf("unexpected original_url: %q", res.OriginalURL)
	}
}

func TestLinkHandler_CreateLink_ProjectNotFound(t *testing.T) {
	t.Parallel()
	h, _, _ := newLinkHandler()

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		withPathValues(
			testutil.NewRequest(http.MethodPost, "/projects/missing/links", dto.LinkRequest{
				URL: "https://example.com",
			}),
			map[string]string{"slug": "missing"},
		),
		int64(1),
	)

	h.CreateLink(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestLinkHandler_CreateLink_ForbiddenOtherUser(t *testing.T) {
	t.Parallel()
	h, _, projects := newLinkHandler()
	seedProject(t, projects, 10, "Owner Project", "owner-project")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		withPathValues(
			testutil.NewRequest(http.MethodPost, "/projects/owner-project/links", dto.LinkRequest{
				URL: "https://example.com",
			}),
			map[string]string{"slug": "owner-project"},
		),
		int64(99),
	)

	h.CreateLink(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

// --- ListLinks ---

func TestLinkHandler_ListLinks_Success(t *testing.T) {
	t.Parallel()
	h, links, projects := newLinkHandler()
	const userID = int64(1)
	p := seedProject(t, projects, userID, "List Project", "list-project")
	seedLink(t, links, p.ID, "https://a.com")
	seedLink(t, links, p.ID, "https://b.com")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		withPathValues(
			testutil.NewRequest(http.MethodGet, "/projects/list-project/links", nil),
			map[string]string{"slug": "list-project"},
		),
		userID,
	)

	h.ListLinks(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var res dto.ListLinksResponse
	if err := testutil.DecodeJSON(w, &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(res.Data) != 2 {
		t.Errorf("expected 2 links, got %d", len(res.Data))
	}
}

// --- GetLink ---

func TestLinkHandler_GetLink_Success(t *testing.T) {
	t.Parallel()
	h, links, projects := newLinkHandler()
	const userID = int64(1)
	p := seedProject(t, projects, userID, "Get Project", "get-project")
	link := seedLink(t, links, p.ID, "https://example.com")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		withPathValues(
			testutil.NewRequest(http.MethodGet, "/projects/get-project/links/"+link.ShortCode, nil),
			map[string]string{"slug": "get-project", "code": link.ShortCode},
		),
		userID,
	)

	h.GetLink(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLinkHandler_GetLink_NotFound(t *testing.T) {
	t.Parallel()
	h, _, projects := newLinkHandler()
	const userID = int64(1)
	seedProject(t, projects, userID, "NF Project", "nf-project")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		withPathValues(
			testutil.NewRequest(http.MethodGet, "/projects/nf-project/links/missing-code", nil),
			map[string]string{"slug": "nf-project", "code": "missing-code"},
		),
		userID,
	)

	h.GetLink(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// --- UpdateLink ---

func TestLinkHandler_UpdateLink_Success(t *testing.T) {
	t.Parallel()
	h, links, projects := newLinkHandler()
	const userID = int64(1)
	p := seedProject(t, projects, userID, "Update Project", "update-project")
	link := seedLink(t, links, p.ID, "https://old.com")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		withPathValues(
			testutil.NewRequest(http.MethodPut, "/projects/update-project/links/"+link.ShortCode, dto.LinkRequest{
				URL: "https://new.com",
			}),
			map[string]string{"slug": "update-project", "code": link.ShortCode},
		),
		userID,
	)

	h.UpdateLink(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var res dto.LinkResponse
	if err := testutil.DecodeJSON(w, &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if res.OriginalURL != "https://new.com" {
		t.Errorf("expected updated URL, got %q", res.OriginalURL)
	}
}

// --- DeleteLink ---

func TestLinkHandler_DeleteLink_Success(t *testing.T) {
	t.Parallel()
	h, links, projects := newLinkHandler()
	const userID = int64(1)
	p := seedProject(t, projects, userID, "Delete Project", "delete-project")
	link := seedLink(t, links, p.ID, "https://todelete.com")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		withPathValues(
			testutil.NewRequest(http.MethodDelete, "/projects/delete-project/links/"+link.ShortCode, nil),
			map[string]string{"slug": "delete-project", "code": link.ShortCode},
		),
		userID,
	)

	h.DeleteLink(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d: %s", w.Code, w.Body.String())
	}

	if _, err := links.GetByCode(t.Context(), link.ShortCode); err == nil {
		t.Error("expected link to be deleted from repo")
	}
}

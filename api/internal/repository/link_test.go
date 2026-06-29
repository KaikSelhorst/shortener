//go:build integration

package repository_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
)

var projectSeq atomic.Int64

func newLinkRepo() *repository.LinkRepository {
	return repository.NewLinkRepository(testDB)
}

func makeProject(t *testing.T, userID int64) *model.Project {
	t.Helper()
	slug := fmt.Sprintf("%s-%d", t.Name(), projectSeq.Add(1))
	p := &model.Project{UserID: userID, Name: slug, Slug: slug}
	if err := newProjectRepo().Create(t.Context(), p); err != nil {
		t.Fatalf("makeProject: %v", err)
	}
	return p
}

func TestLinkRepository_Create(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)

	svc, _ := service.NewShortcodeService([]byte("test-shortcode-secret-32-chars!!"))
	repo := newLinkRepo()

	link := &model.Link{ProjectID: p.ID, OriginalURL: "https://example.com"}
	if err := repo.Create(t.Context(), link, svc.GenerateShortCode); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if link.ID == 0 {
		t.Error("expected ID after create")
	}
	if link.ShortCode == "" {
		t.Error("expected ShortCode to be generated")
	}
}

func TestLinkRepository_GetByCode(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)

	svc, _ := service.NewShortcodeService([]byte("test-shortcode-secret-32-chars!!"))
	repo := newLinkRepo()

	link := &model.Link{ProjectID: p.ID, OriginalURL: "https://getbycode.com"}
	if err := repo.Create(t.Context(), link, svc.GenerateShortCode); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	found, err := repo.GetByCode(t.Context(), link.ShortCode)
	if err != nil {
		t.Fatalf("GetByCode: %v", err)
	}
	if found.OriginalURL != link.OriginalURL {
		t.Errorf("expected URL %q, got %q", link.OriginalURL, found.OriginalURL)
	}
}

func TestLinkRepository_GetByCode_NotFound(t *testing.T) {
	truncate(t)

	_, err := newLinkRepo().GetByCode(t.Context(), "nonexistent-code")
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestLinkRepository_Update(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)

	svc, _ := service.NewShortcodeService([]byte("test-shortcode-secret-32-chars!!"))
	repo := newLinkRepo()

	link := &model.Link{ProjectID: p.ID, OriginalURL: "https://original.com"}
	if err := repo.Create(t.Context(), link, svc.GenerateShortCode); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	link.OriginalURL = "https://updated.com"
	if err := repo.Update(t.Context(), link); err != nil {
		t.Fatalf("Update: %v", err)
	}

	found, err := repo.GetByCode(t.Context(), link.ShortCode)
	if err != nil {
		t.Fatalf("GetByCode after update: %v", err)
	}
	if found.OriginalURL != "https://updated.com" {
		t.Errorf("expected updated URL, got %q", found.OriginalURL)
	}
}

func TestLinkRepository_DeleteByCode(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)

	svc, _ := service.NewShortcodeService([]byte("test-shortcode-secret-32-chars!!"))
	repo := newLinkRepo()

	link := &model.Link{ProjectID: p.ID, OriginalURL: "https://todelete.com"}
	if err := repo.Create(t.Context(), link, svc.GenerateShortCode); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	if err := repo.DeleteByCode(t.Context(), link.ShortCode); err != nil {
		t.Fatalf("DeleteByCode: %v", err)
	}

	_, err := repo.GetByCode(t.Context(), link.ShortCode)
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestLinkRepository_List(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	p := makeProject(t, u.ID)

	svc, _ := service.NewShortcodeService([]byte("test-shortcode-secret-32-chars!!"))
	repo := newLinkRepo()

	for _, url := range []string{"https://a.com", "https://b.com", "https://c.com"} {
		if err := repo.Create(t.Context(), &model.Link{ProjectID: p.ID, OriginalURL: url}, svc.GenerateShortCode); err != nil {
			t.Fatalf("setup Create %q: %v", url, err)
		}
	}

	links, err := repo.List(t.Context(), p.ID, 0, "", 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(links) != 3 {
		t.Errorf("expected 3 links, got %d", len(links))
	}
}

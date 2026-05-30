//go:build integration

package repository_test

import (
	"testing"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
)

func newProjectRepo() *repository.ProjectRepository {
	return repository.NewProjectRepository(testDB)
}

func makeUser(t *testing.T) *model.User {
	t.Helper()
	u := &model.User{Email: t.Name() + "@example.com", PasswordHash: "hash"}
	if err := repository.NewUserRepository(testDB).Create(t.Context(), u); err != nil {
		t.Fatalf("makeUser: %v", err)
	}
	return u
}

func TestProjectRepository_Create(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newProjectRepo()

	p := &model.Project{UserID: u.ID, Name: "My Project", Slug: "my-project"}
	if err := repo.Create(t.Context(), p); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if p.ID == 0 {
		t.Error("expected ID to be set after create")
	}
}

func TestProjectRepository_FindBySlug(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newProjectRepo()

	p := &model.Project{UserID: u.ID, Name: "Slug Project", Slug: "slug-project"}
	if err := repo.Create(t.Context(), p); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	found, err := repo.FindBySlug(t.Context(), "slug-project")
	if err != nil {
		t.Fatalf("FindBySlug: %v", err)
	}
	if found.ID != p.ID {
		t.Errorf("expected ID %d, got %d", p.ID, found.ID)
	}
}

func TestProjectRepository_FindBySlug_NotFound(t *testing.T) {
	truncate(t)

	_, err := newProjectRepo().FindBySlug(t.Context(), "nonexistent")
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestProjectRepository_GetByID(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newProjectRepo()

	p := &model.Project{UserID: u.ID, Name: "ID Project", Slug: "id-project"}
	if err := repo.Create(t.Context(), p); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	found, err := repo.GetByID(t.Context(), p.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if found.Slug != p.Slug {
		t.Errorf("expected slug %q, got %q", p.Slug, found.Slug)
	}
}

func TestProjectRepository_FindAllByUserID(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newProjectRepo()

	for _, slug := range []string{"p1", "p2", "p3"} {
		if err := repo.Create(t.Context(), &model.Project{UserID: u.ID, Name: slug, Slug: slug}); err != nil {
			t.Fatalf("setup Create %q: %v", slug, err)
		}
	}

	projects, err := repo.FindAllByUserID(t.Context(), u.ID)
	if err != nil {
		t.Fatalf("FindAllByUserID: %v", err)
	}
	if len(projects) != 3 {
		t.Errorf("expected 3 projects, got %d", len(projects))
	}
}

func TestProjectRepository_Update(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newProjectRepo()

	p := &model.Project{UserID: u.ID, Name: "Original", Slug: "original"}
	if err := repo.Create(t.Context(), p); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	p.Name = "Updated"
	p.Slug = "updated"
	if err := repo.Update(t.Context(), p); err != nil {
		t.Fatalf("Update: %v", err)
	}

	found, err := repo.FindBySlug(t.Context(), "updated")
	if err != nil {
		t.Fatalf("FindBySlug after update: %v", err)
	}
	if found.Name != "Updated" {
		t.Errorf("expected name Updated, got %q", found.Name)
	}
}

func TestProjectRepository_Delete(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newProjectRepo()

	p := &model.Project{UserID: u.ID, Name: "ToDelete", Slug: "to-delete"}
	if err := repo.Create(t.Context(), p); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	if err := repo.Delete(t.Context(), p.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := repo.GetByID(t.Context(), p.ID)
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

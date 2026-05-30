//go:build integration

package repository_test

import (
	"testing"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
)

func newUserRepo() *repository.UserRepository {
	return repository.NewUserRepository(testDB)
}

func TestUserRepository_Create(t *testing.T) {
	truncate(t)
	repo := newUserRepo()

	user := &model.User{Email: "alice@example.com", PasswordHash: "hash"}
	if err := repo.Create(t.Context(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected ID to be set after create")
	}
	if user.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	truncate(t)
	repo := newUserRepo()

	u := &model.User{Email: "dup@example.com", PasswordHash: "hash"}
	if err := repo.Create(t.Context(), u); err != nil {
		t.Fatalf("first Create: %v", err)
	}

	u2 := &model.User{Email: "dup@example.com", PasswordHash: "hash2"}
	if err := repo.Create(t.Context(), u2); err == nil {
		t.Error("expected error on duplicate email, got nil")
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	truncate(t)
	repo := newUserRepo()

	created := &model.User{Email: "bob@example.com", PasswordHash: "hash"}
	if err := repo.Create(t.Context(), created); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	found, err := repo.FindByEmail(t.Context(), "bob@example.com")
	if err != nil {
		t.Fatalf("FindByEmail: %v", err)
	}
	if found.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, found.ID)
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	truncate(t)

	_, err := newUserRepo().FindByEmail(t.Context(), "nobody@example.com")
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	truncate(t)
	repo := newUserRepo()

	u := &model.User{Email: "carol@example.com", PasswordHash: "hash"}
	if err := repo.Create(t.Context(), u); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	found, err := repo.FindByID(t.Context(), u.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if found.Email != u.Email {
		t.Errorf("expected email %q, got %q", u.Email, found.Email)
	}
}

func TestUserRepository_SaveTOTPSecret(t *testing.T) {
	truncate(t)
	repo := newUserRepo()

	u := &model.User{Email: "dave@example.com", PasswordHash: "hash"}
	if err := repo.Create(t.Context(), u); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	secret := "TESTSECRET123"
	if err := repo.SaveTOTPSecret(t.Context(), u.ID, secret); err != nil {
		t.Fatalf("SaveTOTPSecret: %v", err)
	}

	found, err := repo.FindByID(t.Context(), u.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if found.TOTPSecret == nil || *found.TOTPSecret != secret {
		t.Errorf("expected secret %q, got %v", secret, found.TOTPSecret)
	}
}

func TestUserRepository_SetTOTPEnabled(t *testing.T) {
	truncate(t)
	repo := newUserRepo()

	u := &model.User{Email: "eve@example.com", PasswordHash: "hash"}
	if err := repo.Create(t.Context(), u); err != nil {
		t.Fatalf("setup Create: %v", err)
	}
	if err := repo.SaveTOTPSecret(t.Context(), u.ID, "MYSECRET"); err != nil {
		t.Fatalf("setup SaveTOTPSecret: %v", err)
	}

	if err := repo.SetTOTPEnabled(t.Context(), u.ID, true); err != nil {
		t.Fatalf("SetTOTPEnabled(true): %v", err)
	}
	found, err := repo.FindByID(t.Context(), u.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if !found.TOTPEnabled {
		t.Error("expected TOTPEnabled=true")
	}

	if err := repo.SetTOTPEnabled(t.Context(), u.ID, false); err != nil {
		t.Fatalf("SetTOTPEnabled(false): %v", err)
	}
	found, err = repo.FindByID(t.Context(), u.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if found.TOTPEnabled {
		t.Error("expected TOTPEnabled=false after disable")
	}
	if found.TOTPSecret != nil {
		t.Error("expected TOTPSecret=nil after disable")
	}
}

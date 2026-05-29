//go:build integration

package repository_test

import (
	"testing"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
)

func newAPIKeyRepo() *repository.APIKeyRepository {
	return repository.NewAPIKeyRepository(testDB)
}

func makeAPIKey(t *testing.T, userID int64) *model.APIKey {
	t.Helper()
	k := &model.APIKey{
		UserID:    userID,
		Name:      t.Name(),
		KeyPrefix: "sk_test",
		KeyHash:   t.Name() + "-hash",
		Scopes:    []string{"links:read", "links:write"},
	}
	if err := newAPIKeyRepo().Create(t.Context(), k); err != nil {
		t.Fatalf("makeAPIKey: %v", err)
	}
	return k
}

func TestAPIKeyRepository_Create(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newAPIKeyRepo()

	k := &model.APIKey{
		UserID:    u.ID,
		Name:      "My Key",
		KeyPrefix: "sk_abcd",
		KeyHash:   "sha256hashvalue",
		Scopes:    []string{"links:read"},
	}
	if err := repo.Create(t.Context(), k); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if k.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestAPIKeyRepository_GetByHash(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	k := makeAPIKey(t, u.ID)

	found, err := newAPIKeyRepo().GetByHash(t.Context(), k.KeyHash)
	if err != nil {
		t.Fatalf("GetByHash: %v", err)
	}
	if found.Name != k.Name {
		t.Errorf("expected name %q, got %q", k.Name, found.Name)
	}
}

func TestAPIKeyRepository_GetByHash_NotFound(t *testing.T) {
	truncate(t)

	_, err := newAPIKeyRepo().GetByHash(t.Context(), "nonexistent-hash")
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAPIKeyRepository_List(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newAPIKeyRepo()

	for i := 0; i < 3; i++ {
		k := &model.APIKey{
			UserID:    u.ID,
			Name:      "key",
			KeyPrefix: "sk_x",
			KeyHash:   "hash" + string(rune('0'+i)),
			Scopes:    []string{"*"},
		}
		if err := repo.Create(t.Context(), k); err != nil {
			t.Fatalf("setup Create key %d: %v", i, err)
		}
	}

	keys, err := repo.List(t.Context(), u.ID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}

func TestAPIKeyRepository_Delete(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	k := makeAPIKey(t, u.ID)
	repo := newAPIKeyRepo()

	if err := repo.Delete(t.Context(), k.ID, u.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := repo.GetByHash(t.Context(), k.KeyHash)
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestAPIKeyRepository_Delete_WrongUser(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	k := makeAPIKey(t, u.ID)

	err := newAPIKeyRepo().Delete(t.Context(), k.ID, u.ID+999)
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound for wrong user, got %v", err)
	}
}

func TestAPIKeyRepository_UpdateLastUsed(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	k := makeAPIKey(t, u.ID)
	repo := newAPIKeyRepo()

	if err := repo.UpdateLastUsed(t.Context(), k.ID); err != nil {
		t.Fatalf("UpdateLastUsed: %v", err)
	}

	found, err := repo.GetByHash(t.Context(), k.KeyHash)
	if err != nil {
		t.Fatalf("GetByHash: %v", err)
	}
	if found.LastUsedAt == nil {
		t.Error("expected LastUsedAt to be set after UpdateLastUsed")
	}
}

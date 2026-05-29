//go:build integration

package repository_test

import (
	"testing"
	"time"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
)

func newRefreshTokenRepo() *repository.RefreshTokenRepository {
	return repository.NewRefreshTokenRepository(testDB)
}

func TestRefreshTokenRepository_Create(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newRefreshTokenRepo()

	rt := &model.RefreshToken{
		UserID:    u.ID,
		TokenHash: "abc123hash",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := repo.Create(t.Context(), rt); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if rt.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestRefreshTokenRepository_RevokeIfActive(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newRefreshTokenRepo()

	rt := &model.RefreshToken{
		UserID:    u.ID,
		TokenHash: "revoke-me-hash",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if err := repo.Create(t.Context(), rt); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	revoked, err := repo.RevokeIfActive(t.Context(), "revoke-me-hash")
	if err != nil {
		t.Fatalf("RevokeIfActive: %v", err)
	}
	if revoked.UserID != u.ID {
		t.Errorf("expected userID %d, got %d", u.ID, revoked.UserID)
	}
}

func TestRefreshTokenRepository_RevokeIfActive_AlreadyRevoked(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newRefreshTokenRepo()

	rt := &model.RefreshToken{
		UserID:    u.ID,
		TokenHash: "already-revoked",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if err := repo.Create(t.Context(), rt); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	if _, err := repo.RevokeIfActive(t.Context(), "already-revoked"); err != nil {
		t.Fatalf("first revoke: %v", err)
	}

	_, err := repo.RevokeIfActive(t.Context(), "already-revoked")
	if err == nil {
		t.Error("expected error on second revocation, got nil")
	}
}

func TestRefreshTokenRepository_RevokeIfActive_Expired(t *testing.T) {
	truncate(t)
	u := makeUser(t)
	repo := newRefreshTokenRepo()

	rt := &model.RefreshToken{
		UserID:    u.ID,
		TokenHash: "expired-token",
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	if err := repo.Create(t.Context(), rt); err != nil {
		t.Fatalf("setup Create: %v", err)
	}

	_, err := repo.RevokeIfActive(t.Context(), "expired-token")
	if err == nil {
		t.Error("expected error revoking expired token, got nil")
	}
}

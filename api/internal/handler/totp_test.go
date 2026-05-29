package handler_test

import (
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

func newTOTPHandler() (*handler.TOTPHandler, *fakes.UserRepo, *fakes.RefreshTokenRepo) {
	users := fakes.NewUserRepo()
	tokens := fakes.NewRefreshTokenRepo()
	svc := service.NewAuthService("test-secret")
	return handler.NewTOTPHandler(users, tokens, svc), users, tokens
}

func enableTOTP(t *testing.T, users *fakes.UserRepo, userID int64, secret string) {
	t.Helper()
	if err := users.SaveTOTPSecret(t.Context(), userID, secret); err != nil {
		t.Fatalf("enableTOTP: SaveTOTPSecret: %v", err)
	}
	if err := users.SetTOTPEnabled(t.Context(), userID, true); err != nil {
		t.Fatalf("enableTOTP: SetTOTPEnabled: %v", err)
	}
}

// --- Setup ---

func TestTOTPHandler_Setup_Success(t *testing.T) {
	t.Parallel()
	h, users, _ := newTOTPHandler()
	u := seedUser(t, users, "alice@totp.com", "password")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		testutil.NewRequest(http.MethodPost, "/auth/totp/setup", nil),
		u.ID,
	)

	h.Setup(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var res dto.TOTPSetupResponse
	if err := testutil.DecodeJSON(w, &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if res.URI == "" {
		t.Error("expected URI in setup response")
	}
	if res.Secret == "" {
		t.Error("expected secret in setup response")
	}
}

func TestTOTPHandler_Setup_BlockedByAPIKey(t *testing.T) {
	t.Parallel()
	h, users, _ := newTOTPHandler()
	u := seedUser(t, users, "bob@totp.com", "password")

	w := httptest.NewRecorder()
	r := testutil.WithAPIKey(
		testutil.NewRequest(http.MethodPost, "/auth/totp/setup", nil),
		&model.APIKey{UserID: u.ID, Scopes: []string{"*"}},
	)

	h.Setup(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestTOTPHandler_Setup_AlreadyEnabled(t *testing.T) {
	t.Parallel()
	h, users, _ := newTOTPHandler()
	u := seedUser(t, users, "carol@totp.com", "password")
	enableTOTP(t, users, u.ID, "JBSWY3DPEHPK3PXP")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		testutil.NewRequest(http.MethodPost, "/auth/totp/setup", nil),
		u.ID,
	)

	h.Setup(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

// --- Confirm ---

func TestTOTPHandler_Confirm_Success(t *testing.T) {
	t.Parallel()
	h, users, _ := newTOTPHandler()
	u := seedUser(t, users, "dave@totp.com", "password")

	secret, _ := service.GenerateTOTPSecret()
	if err := users.SaveTOTPSecret(t.Context(), u.ID, secret); err != nil {
		t.Fatalf("SaveTOTPSecret: %v", err)
	}

	code, err := testutil.GenerateCurrentTOTPCode(secret)
	if err != nil {
		t.Fatalf("generate code: %v", err)
	}

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		testutil.NewRequest(http.MethodPost, "/auth/totp/confirm", dto.TOTPConfirmRequest{Code: code}),
		u.ID,
	)

	h.Confirm(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d: %s", w.Code, w.Body.String())
	}

	updated, err := users.FindByID(t.Context(), u.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if !updated.TOTPEnabled {
		t.Error("expected TOTPEnabled=true after confirm")
	}
}

func TestTOTPHandler_Confirm_NoSetup(t *testing.T) {
	t.Parallel()
	h, users, _ := newTOTPHandler()
	u := seedUser(t, users, "eve@totp.com", "password")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		testutil.NewRequest(http.MethodPost, "/auth/totp/confirm", dto.TOTPConfirmRequest{Code: "123456"}),
		u.ID,
	)

	h.Confirm(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 when setup not initiated, got %d", w.Code)
	}
}

func TestTOTPHandler_Confirm_AlreadyEnabled(t *testing.T) {
	t.Parallel()
	h, users, _ := newTOTPHandler()
	u := seedUser(t, users, "frank@totp.com", "password")
	enableTOTP(t, users, u.ID, "JBSWY3DPEHPK3PXP")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		testutil.NewRequest(http.MethodPost, "/auth/totp/confirm", dto.TOTPConfirmRequest{Code: "000000"}),
		u.ID,
	)

	h.Confirm(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 when already enabled, got %d", w.Code)
	}
}

// --- Disable ---

func TestTOTPHandler_Disable_Success(t *testing.T) {
	t.Parallel()
	h, users, _ := newTOTPHandler()
	u := seedUser(t, users, "grace@totp.com", "password")

	secret, _ := service.GenerateTOTPSecret()
	enableTOTP(t, users, u.ID, secret)

	code, err := testutil.GenerateCurrentTOTPCode(secret)
	if err != nil {
		t.Fatalf("generate code: %v", err)
	}

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		testutil.NewRequest(http.MethodDelete, "/auth/totp", dto.TOTPDisableRequest{Code: code}),
		u.ID,
	)

	h.Disable(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d: %s", w.Code, w.Body.String())
	}

	updated, err := users.FindByID(t.Context(), u.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if updated.TOTPEnabled {
		t.Error("expected TOTPEnabled=false after disable")
	}
}

func TestTOTPHandler_Disable_NotEnabled(t *testing.T) {
	t.Parallel()
	h, users, _ := newTOTPHandler()
	u := seedUser(t, users, "henry@totp.com", "password")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		testutil.NewRequest(http.MethodDelete, "/auth/totp", dto.TOTPDisableRequest{Code: "123456"}),
		u.ID,
	)

	h.Disable(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 when TOTP not enabled, got %d", w.Code)
	}
}

func TestTOTPHandler_Disable_BlockedByAPIKey(t *testing.T) {
	t.Parallel()
	h, users, _ := newTOTPHandler()
	u := seedUser(t, users, "ivan@totp.com", "password")

	w := httptest.NewRecorder()
	r := testutil.WithAPIKey(
		testutil.NewRequest(http.MethodDelete, "/auth/totp", dto.TOTPDisableRequest{Code: "123456"}),
		&model.APIKey{UserID: u.ID, Scopes: []string{"*"}},
	)

	h.Disable(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

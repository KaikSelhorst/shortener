package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/handler"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository/fakes"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/KaikSelhorst/shortener/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

func newAuthHandler() (*handler.AuthHandler, *fakes.UserRepo, *fakes.RefreshTokenRepo) {
	users := fakes.NewUserRepo()
	tokens := fakes.NewRefreshTokenRepo()
	svc := service.NewAuthService("test-secret")
	return handler.NewAuthHandler(users, tokens, svc), users, tokens
}

func seedUser(t *testing.T, users *fakes.UserRepo, email, password string) *model.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	u := &model.User{Email: email, PasswordHash: string(hash)}
	if err := users.Create(t.Context(), u); err != nil {
		t.Fatalf("seedUser: %v", err)
	}
	return u
}

// --- Register ---

func TestAuthHandler_Register_Success(t *testing.T) {
	h, _, _ := newAuthHandler()
	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/register", dto.RegisterRequest{
		Email:    "alice@example.com",
		Password: "supersecret",
	})

	h.Register(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var res dto.AuthState
	if err := testutil.DecodeJSON(w, &res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res.Next != "complete" {
		t.Errorf("expected next=complete, got %q", res.Next)
	}
	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Error("expected tokens to be present after registration")
	}
}

func TestAuthHandler_Register_InvalidPayload(t *testing.T) {
	h, _, _ := newAuthHandler()
	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/register", dto.RegisterRequest{
		Email:    "not-an-email",
		Password: "short",
	})

	h.Register(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_Register_DuplicateEmail(t *testing.T) {
	h, users, _ := newAuthHandler()
	seedUser(t, users, "alice@example.com", "password123")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/register", dto.RegisterRequest{
		Email:    "alice@example.com",
		Password: "password123",
	})

	h.Register(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Login ---

func TestAuthHandler_Login_NoTOTP(t *testing.T) {
	h, users, _ := newAuthHandler()
	seedUser(t, users, "bob@example.com", "mypassword")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/login", dto.LoginRequest{
		Email:    "bob@example.com",
		Password: "mypassword",
	})

	h.Login(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var res dto.AuthState
	_ = testutil.DecodeJSON(w, &res)
	if res.Next != "complete" {
		t.Errorf("expected next=complete, got %q", res.Next)
	}
}

func TestAuthHandler_Login_WithTOTP(t *testing.T) {
	h, users, _ := newAuthHandler()
	u := seedUser(t, users, "carol@example.com", "mypassword")

	secret := "JBSWY3DPEHPK3PXP"
	_ = users.SaveTOTPSecret(nil, u.ID, secret)
	_ = users.SetTOTPEnabled(nil, u.ID, true)

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/login", dto.LoginRequest{
		Email:    "carol@example.com",
		Password: "mypassword",
	})

	h.Login(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var res dto.AuthState
	_ = testutil.DecodeJSON(w, &res)
	if res.Next != "totp" {
		t.Errorf("expected next=totp, got %q", res.Next)
	}
	if res.Session == "" {
		t.Error("expected session token to be present")
	}
	if res.AccessToken != "" {
		t.Error("access token should not be issued before TOTP")
	}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	h, users, _ := newAuthHandler()
	seedUser(t, users, "dave@example.com", "correctpassword")

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/login", dto.LoginRequest{
		Email:    "dave@example.com",
		Password: "wrongpassword",
	})

	h.Login(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_Login_UnknownEmail(t *testing.T) {
	h, _, _ := newAuthHandler()

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/login", dto.LoginRequest{
		Email:    "nobody@example.com",
		Password: "anypassword",
	})

	h.Login(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// --- Refresh ---

func TestAuthHandler_Refresh_ValidToken(t *testing.T) {
	h, users, _ := newAuthHandler()
	seedUser(t, users, "eve@example.com", "password123")

	// Login to get a refresh token.
	loginW := httptest.NewRecorder()
	h.Login(loginW, testutil.NewRequest(http.MethodPost, "/auth/login", dto.LoginRequest{
		Email:    "eve@example.com",
		Password: "password123",
	}))

	var loginRes dto.AuthState
	_ = json.NewDecoder(loginW.Body).Decode(&loginRes)

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/refresh", dto.RefreshRequest{
		RefreshToken: loginRes.RefreshToken,
	})

	h.Refresh(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var res dto.AuthState
	_ = testutil.DecodeJSON(w, &res)
	if res.AccessToken == "" {
		t.Error("expected new access token after refresh")
	}
}

func TestAuthHandler_Refresh_InvalidToken(t *testing.T) {
	h, _, _ := newAuthHandler()

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/refresh", dto.RefreshRequest{
		RefreshToken: "invalid-token",
	})

	h.Refresh(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// --- Me ---

func TestAuthHandler_Me_WithJWT(t *testing.T) {
	h, users, _ := newAuthHandler()
	u := seedUser(t, users, "frank@example.com", "password")

	w := httptest.NewRecorder()
	r := testutil.WithUserID(
		testutil.NewRequest(http.MethodGet, "/auth/me", nil),
		u.ID,
	)

	h.Me(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_Me_WithAPIKey(t *testing.T) {
	h, users, _ := newAuthHandler()
	u := seedUser(t, users, "grace@example.com", "password")

	apiKey := &model.APIKey{UserID: u.ID, Scopes: []string{"*"}}

	w := httptest.NewRecorder()
	r := testutil.WithAPIKey(
		testutil.NewRequest(http.MethodGet, "/auth/me", nil),
		apiKey,
	)

	h.Me(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for API key on /me, got %d", w.Code)
	}
}

// --- Logout ---

func TestAuthHandler_Logout(t *testing.T) {
	h, users, _ := newAuthHandler()
	seedUser(t, users, "heidi@example.com", "password123")

	loginW := httptest.NewRecorder()
	h.Login(loginW, testutil.NewRequest(http.MethodPost, "/auth/login", dto.LoginRequest{
		Email:    "heidi@example.com",
		Password: "password123",
	}))

	var loginRes dto.AuthState
	_ = json.NewDecoder(loginW.Body).Decode(&loginRes)

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/logout", map[string]string{
		"refresh_token": loginRes.RefreshToken,
	})
	r.Header.Set("Authorization", "Bearer "+loginRes.AccessToken)

	h.Logout(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestAuthHandler_Logout_RevokesAccessToken(t *testing.T) {
	svc := service.NewAuthService("test-secret")
	users := fakes.NewUserRepo()
	tokens := fakes.NewRefreshTokenRepo()
	h := handler.NewAuthHandler(users, tokens, svc)

	u := seedUser(t, users, "ivan@example.com", "password123")

	accessToken, _ := svc.GenerateAccessToken(u.ID)
	refreshRaw, refreshHash, _ := svc.GenerateRefreshToken()
	if err := tokens.Create(t.Context(), &model.RefreshToken{
		UserID:    u.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(service.RefreshTokenTTL),
	}); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if _, err := svc.ValidateAccessToken(accessToken); err != nil {
		t.Fatalf("expected valid before logout: %v", err)
	}

	w := httptest.NewRecorder()
	r := testutil.NewRequest(http.MethodPost, "/auth/logout", map[string]string{
		"refresh_token": refreshRaw,
	})
	r.Header.Set("Authorization", "Bearer "+accessToken)
	h.Logout(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	if _, err := svc.ValidateAccessToken(accessToken); err == nil {
		t.Error("expected access token to be invalid after logout")
	}
}

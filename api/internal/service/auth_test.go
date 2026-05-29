package service

import (
	"strings"
	"testing"
	"time"
)

func newTestAuthService() *AuthService {
	return NewAuthService("test-jwt-secret-key")
}

// --- Access token ---

func TestAuthService_AccessToken_RoundTrip(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	userID := int64(42)

	token, err := svc.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	got, err := svc.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken: %v", err)
	}
	if got != userID {
		t.Errorf("expected userID %d, got %d", userID, got)
	}
}

func TestAuthService_AccessToken_WrongSecret(t *testing.T) {
	t.Parallel()
	svc1 := NewAuthService("secret-a")
	svc2 := NewAuthService("secret-b")

	token, _ := svc1.GenerateAccessToken(1)
	_, err := svc2.ValidateAccessToken(token)
	if err == nil {
		t.Error("expected error: token signed with different secret should be rejected")
	}
}

func TestAuthService_AccessToken_RejectsSessionToken(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	sessionToken, _ := svc.GenerateSessionToken(1, "totp")

	_, err := svc.ValidateAccessToken(sessionToken)
	if err == nil {
		t.Error("expected error: session token must not be accepted as access token")
	}
}

// --- Session token ---

func TestAuthService_SessionToken_RoundTrip(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	token, err := svc.GenerateSessionToken(99, "totp")
	if err != nil {
		t.Fatalf("GenerateSessionToken: %v", err)
	}

	userID, next, err := svc.ValidateSessionToken(token)
	if err != nil {
		t.Fatalf("ValidateSessionToken: %v", err)
	}
	if userID != 99 {
		t.Errorf("expected userID 99, got %d", userID)
	}
	if next != "totp" {
		t.Errorf("expected next=totp, got %q", next)
	}
}

func TestAuthService_SessionToken_SingleUse(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	token, _ := svc.GenerateSessionToken(99, "totp")

	if _, _, err := svc.ValidateSessionToken(token); err != nil {
		t.Fatalf("first use: %v", err)
	}
	if _, _, err := svc.ValidateSessionToken(token); err == nil {
		t.Error("expected error on second use of the same session token")
	}
}

func TestAuthService_ParseSessionToken_DoesNotConsume(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	token, _ := svc.GenerateSessionToken(99, "totp")

	// Parsing twice must not mark the token as consumed.
	for range 2 {
		userID, next, jti, expiry, err := svc.ParseSessionToken(token)
		if err != nil {
			t.Fatalf("ParseSessionToken: %v", err)
		}
		if userID != 99 || next != "totp" || jti == "" || expiry.IsZero() {
			t.Errorf("unexpected values: userID=%d next=%q jti=%q expiry=%v", userID, next, jti, expiry)
		}
	}
}

func TestAuthService_ConsumeSession_SingleUse(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	token, _ := svc.GenerateSessionToken(99, "totp")

	_, _, jti, expiry, err := svc.ParseSessionToken(token)
	if err != nil {
		t.Fatalf("ParseSessionToken: %v", err)
	}

	if err := svc.ConsumeSession(jti, expiry); err != nil {
		t.Fatalf("first ConsumeSession: %v", err)
	}
	if err := svc.ConsumeSession(jti, expiry); err == nil {
		t.Error("expected error on second ConsumeSession with same JTI")
	}
}

func TestAuthService_RevokeAccessToken(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()

	tokenStr, _ := svc.GenerateAccessToken(42)

	// Token is valid before revocation.
	if _, err := svc.ValidateAccessToken(tokenStr); err != nil {
		t.Fatalf("expected valid before revoke: %v", err)
	}

	if err := svc.RevokeAccessToken(tokenStr); err != nil {
		t.Fatalf("RevokeAccessToken: %v", err)
	}

	// Token must be rejected after revocation.
	if _, err := svc.ValidateAccessToken(tokenStr); err == nil {
		t.Error("expected error after revocation, got nil")
	}
}

func TestAuthService_SessionToken_RejectsAccessToken(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	accessToken, _ := svc.GenerateAccessToken(1)

	_, _, err := svc.ValidateSessionToken(accessToken)
	if err == nil {
		t.Error("expected error: access token must not be accepted as session token")
	}
}

// --- Refresh token ---

func TestAuthService_GenerateRefreshToken_Format(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	raw, hash, err := svc.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}
	if len(raw) != 64 {
		t.Errorf("expected raw length 64, got %d", len(raw))
	}
	if hash != HashToken(raw) {
		t.Error("returned hash does not match HashToken(raw)")
	}
}

func TestAuthService_GenerateRefreshToken_Unique(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	raw1, _, _ := svc.GenerateRefreshToken()
	raw2, _, _ := svc.GenerateRefreshToken()
	if raw1 == raw2 {
		t.Error("GenerateRefreshToken returned the same token twice")
	}
}

// --- API key ---

func TestAuthService_GenerateAPIKey_Format(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	raw, hash, err := svc.GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey: %v", err)
	}
	if !strings.HasPrefix(raw, "sk_") {
		t.Errorf("expected sk_ prefix, got: %q", raw)
	}
	// "sk_" + 64 hex chars
	if len(raw) != 67 {
		t.Errorf("expected length 67, got %d: %q", len(raw), raw)
	}
	if hash != HashToken(raw) {
		t.Error("returned hash does not match HashToken(raw)")
	}
}

func TestAuthService_GenerateAPIKey_Unique(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	raw1, _, _ := svc.GenerateAPIKey()
	raw2, _, _ := svc.GenerateAPIKey()
	if raw1 == raw2 {
		t.Error("GenerateAPIKey returned the same key twice")
	}
}

// --- TOTP replay prevention ---

func TestAuthService_ValidateTOTPAndConsume_Valid(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	secret, _ := GenerateTOTPSecret()
	counter := time.Now().Unix() / totpPeriod
	code, _ := generateTOTP(secret, counter)

	ok, err := svc.ValidateTOTPAndConsume(1, secret, code)
	if err != nil {
		t.Fatalf("ValidateTOTPAndConsume: %v", err)
	}
	if !ok {
		t.Error("expected valid TOTP code to be accepted")
	}
}

func TestAuthService_ValidateTOTPAndConsume_ReplayPrevention(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	secret, _ := GenerateTOTPSecret()
	counter := time.Now().Unix() / totpPeriod
	code, _ := generateTOTP(secret, counter)

	ok, _ := svc.ValidateTOTPAndConsume(1, secret, code)
	if !ok {
		t.Fatal("first call should succeed")
	}

	ok, err := svc.ValidateTOTPAndConsume(1, secret, code)
	if err != nil {
		t.Fatalf("unexpected error on replay: %v", err)
	}
	if ok {
		t.Error("expected replay of same TOTP code to be rejected")
	}
}

func TestAuthService_ValidateTOTPAndConsume_SameCodeDifferentUsers(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	// User 1 consumes their code.
	secret1, _ := GenerateTOTPSecret()
	counter := time.Now().Unix() / totpPeriod
	code1, _ := generateTOTP(secret1, counter)
	ok, _ := svc.ValidateTOTPAndConsume(1, secret1, code1)
	if !ok {
		t.Fatal("user 1: expected success")
	}

	// User 2 with a different secret (different code) should not be affected.
	secret2, _ := GenerateTOTPSecret()
	code2, _ := generateTOTP(secret2, counter)
	ok2, _ := svc.ValidateTOTPAndConsume(2, secret2, code2)
	if !ok2 {
		t.Error("user 2 should be allowed independently")
	}
}

func TestAuthService_ValidateTOTPAndConsume_InvalidCode(t *testing.T) {
	t.Parallel()
	svc := newTestAuthService()
	secret, _ := GenerateTOTPSecret()
	// Code from 2 periods ago is outside the acceptance window.
	counter := time.Now().Unix()/totpPeriod - 2
	staleCode, _ := generateTOTP(secret, counter)

	ok, err := svc.ValidateTOTPAndConsume(1, secret, staleCode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected stale code to be rejected")
	}
}

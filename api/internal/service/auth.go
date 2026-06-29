package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/KaikSelhorst/shortener/internal/jwt"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
	SessionTokenTTL = 5 * time.Minute

	purposeAccess  = "access"
	purposeSession = "mfa_session"
)

// AuthService handles JWT generation/validation and single-use enforcement for
// both MFA session tokens (JTI) and TOTP codes (userID+code).
type AuthService struct {
	jwtSecret     []byte
	usedSessions  sync.Map // jti (string) → expiry (time.Time)
	revokedTokens sync.Map // jti (string) → expiry (time.Time) for revoked access tokens
	usedTOTPCodes sync.Map // "userID:code" (string) → expiry (time.Time)
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{jwtSecret: []byte(jwtSecret)}
}

// accessClaims is the payload for regular access tokens.
// JTI enables individual token revocation (e.g. on logout).
// Purpose prevents session tokens from being accepted as access tokens.
type accessClaims struct {
	UserID  int64              `json:"user_id"`
	Purpose string             `json:"purpose"`
	JTI     string             `json:"jti"`
	Iat     jwt.NumericDate    `json:"iat"`
	Exp     jwt.NumericDate    `json:"exp"`
}

func (s *AuthService) GenerateAccessToken(userID int64) (string, error) {
	jti, err := randomHex(16)
	if err != nil {
		return "", fmt.Errorf("generate access token jti: %w", err)
	}
	now := time.Now()
	return jwt.Sign(accessClaims{
		UserID:  userID,
		Purpose: purposeAccess,
		JTI:     jti,
		Iat:     jwt.NewNumericDate(now),
		Exp:     jwt.NewNumericDate(now.Add(AccessTokenTTL)),
	}, s.jwtSecret)
}

func (s *AuthService) ValidateAccessToken(tokenStr string) (int64, error) {
	var c accessClaims
	if err := jwt.Verify(tokenStr, s.jwtSecret, &c); err != nil {
		return 0, err
	}
	if c.Purpose != purposeAccess {
		return 0, fmt.Errorf("invalid token purpose")
	}
	if time.Now().After(c.Exp.Time()) {
		return 0, fmt.Errorf("token expired")
	}
	if _, revoked := s.revokedTokens.Load(c.JTI); revoked {
		return 0, fmt.Errorf("token revoked")
	}
	return c.UserID, nil
}

// RevokeAccessToken invalidates an access token before its natural TTL expires.
// Called during logout so the token cannot be reused within the 15-minute window.
// Silently ignores expired or malformed tokens since they are already invalid.
func (s *AuthService) RevokeAccessToken(tokenStr string) error {
	var c accessClaims
	if err := jwt.Verify(tokenStr, s.jwtSecret, &c); err != nil {
		return nil
	}
	if c.Purpose != purposeAccess || c.JTI == "" {
		return nil
	}
	s.revokedTokens.Store(c.JTI, c.Exp.Time())
	s.evictRevokedTokens()
	return nil
}

func (s *AuthService) evictRevokedTokens() {
	now := time.Now()
	s.revokedTokens.Range(func(key, value any) bool {
		if exp, ok := value.(time.Time); ok && now.After(exp) {
			s.revokedTokens.Delete(key)
		}
		return true
	})
}

func (s *AuthService) GenerateRefreshToken() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}
	raw = hex.EncodeToString(b)
	hash = HashToken(raw)
	return raw, hash, nil
}

// GenerateAPIKey generates an API Key token in the format sk_<64-chars-hex>.
// Returns the raw token (for one-time display) and its SHA256 hash (for storage).
func (s *AuthService) GenerateAPIKey() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate api key: %w", err)
	}
	raw = "sk_" + hex.EncodeToString(b)
	hash = HashToken(raw)
	return raw, hash, nil
}

func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// sessionClaims holds the intermediate state between authentication steps.
// JTI enforces single-use; Purpose prevents cross-type token confusion.
type sessionClaims struct {
	UserID  int64           `json:"user_id"`
	Next    string          `json:"next"`
	Purpose string          `json:"purpose"`
	JTI     string          `json:"jti"`
	Iat     jwt.NumericDate `json:"iat"`
	Exp     jwt.NumericDate `json:"exp"`
}

// GenerateSessionToken issues a short-lived, single-use JWT that encodes the
// next required authentication step (e.g. "totp").
func (s *AuthService) GenerateSessionToken(userID int64, next string) (string, error) {
	jti, err := randomHex(16)
	if err != nil {
		return "", fmt.Errorf("generate session jti: %w", err)
	}
	now := time.Now()
	return jwt.Sign(sessionClaims{
		UserID:  userID,
		Next:    next,
		Purpose: purposeSession,
		JTI:     jti,
		Iat:     jwt.NewNumericDate(now),
		Exp:     jwt.NewNumericDate(now.Add(SessionTokenTTL)),
	}, s.jwtSecret)
}

// ParseSessionToken parses and validates a session token without consuming its JTI.
// Call ConsumeSession separately to mark the token as used.
// This allows the caller to validate other conditions (e.g. TOTP code) before
// committing the session as consumed, so a wrong TOTP code does not burn the token.
func (s *AuthService) ParseSessionToken(tokenStr string) (userID int64, next string, jti string, expiry time.Time, err error) {
	var c sessionClaims
	if err = jwt.Verify(tokenStr, s.jwtSecret, &c); err != nil {
		return 0, "", "", time.Time{}, err
	}
	if c.Purpose != purposeSession {
		return 0, "", "", time.Time{}, fmt.Errorf("invalid token purpose")
	}
	if time.Now().After(c.Exp.Time()) {
		return 0, "", "", time.Time{}, fmt.Errorf("token expired")
	}
	if _, loaded := s.usedSessions.Load(c.JTI); loaded {
		return 0, "", "", time.Time{}, fmt.Errorf("session token already used")
	}
	return c.UserID, c.Next, c.JTI, c.Exp.Time(), nil
}

// ConsumeSession marks a session JTI as used. Returns an error if the JTI was
// already consumed by a concurrent request.
func (s *AuthService) ConsumeSession(jti string, expiry time.Time) error {
	s.evictExpiredSessions()
	if _, loaded := s.usedSessions.LoadOrStore(jti, expiry); loaded {
		return fmt.Errorf("session token already used")
	}
	return nil
}

// ValidateSessionToken parses, validates, and atomically consumes a session token.
// Use ParseSessionToken + ConsumeSession when you need to validate other conditions
// between the two steps (e.g. the ValidateMFA flow).
func (s *AuthService) ValidateSessionToken(tokenStr string) (int64, string, error) {
	userID, next, jti, expiry, err := s.ParseSessionToken(tokenStr)
	if err != nil {
		return 0, "", err
	}
	if err := s.ConsumeSession(jti, expiry); err != nil {
		return 0, "", err
	}
	return userID, next, nil
}

func (s *AuthService) evictExpiredSessions() {
	now := time.Now()
	s.usedSessions.Range(func(key, value any) bool {
		if exp, ok := value.(time.Time); ok && now.After(exp) {
			s.usedSessions.Delete(key)
		}
		return true
	})
}

// totpReplayWindow is the maximum duration a TOTP code remains valid.
// It covers the ±totpWindow step acceptance window (3 steps × 30 s = 90 s).
const totpReplayWindow = 90 * time.Second

// ValidateTOTPAndConsume validates a TOTP code and enforces single-use within
// the acceptance window. Returns false if the code is invalid or was already
// used by this user. An error is returned only for malformed secrets.
func (s *AuthService) ValidateTOTPAndConsume(userID int64, secret, code string) (bool, error) {
	s.evictUsedTOTPCodes()

	ok, err := ValidateTOTP(secret, code)
	if err != nil || !ok {
		return ok, err
	}

	// Reject replay: store "userID:code" → expiry. LoadOrStore returns loaded=true
	// if another goroutine (or a previous request) already stored this key.
	key := fmt.Sprintf("%d:%s", userID, code)
	expiry := time.Now().Add(totpReplayWindow)
	if _, loaded := s.usedTOTPCodes.LoadOrStore(key, expiry); loaded {
		return false, nil
	}

	return true, nil
}

func (s *AuthService) evictUsedTOTPCodes() {
	now := time.Now()
	s.usedTOTPCodes.Range(func(key, value any) bool {
		if exp, ok := value.(time.Time); ok && now.After(exp) {
			s.usedTOTPCodes.Delete(key)
		}
		return true
	})
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	jwtSecret    []byte
	usedSessions sync.Map // jti (string) → expiry (time.Time)
	usedTOTPCodes sync.Map // "userID:code" (string) → expiry (time.Time)
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{jwtSecret: []byte(jwtSecret)}
}

// claims is the payload for regular access tokens.
// The Purpose field prevents session tokens from being accepted as access tokens.
type claims struct {
	UserID  int64  `json:"user_id"`
	Purpose string `json:"purpose"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateAccessToken(userID int64) (string, error) {
	now := time.Now()
	c := claims{
		UserID:  userID,
		Purpose: purposeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateAccessToken(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Reject session tokens presented as access tokens.
	if c.Purpose != purposeAccess {
		return 0, fmt.Errorf("invalid token purpose")
	}

	return c.UserID, nil
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
	UserID  int64  `json:"user_id"`
	Next    string `json:"next"`
	Purpose string `json:"purpose"`
	JTI     string `json:"jti"`
	jwt.RegisteredClaims
}

// GenerateSessionToken issues a short-lived, single-use JWT that encodes the
// next required authentication step (e.g. "totp").
func (s *AuthService) GenerateSessionToken(userID int64, next string) (string, error) {
	jtiBytes := make([]byte, 16)
	if _, err := rand.Read(jtiBytes); err != nil {
		return "", fmt.Errorf("generate session jti: %w", err)
	}
	jti := hex.EncodeToString(jtiBytes)

	now := time.Now()
	expiry := now.Add(SessionTokenTTL)
	c := sessionClaims{
		UserID:  userID,
		Next:    next,
		Purpose: purposeSession,
		JTI:     jti,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(s.jwtSecret)
}

// ValidateSessionToken parses, validates, and consumes a session token.
// It returns the embedded userID and the next required step.
// Each token may only be used once; a second call with the same token returns an error.
func (s *AuthService) ValidateSessionToken(tokenStr string) (int64, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &sessionClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, "", err
	}

	c, ok := token.Claims.(*sessionClaims)
	if !ok || !token.Valid {
		return 0, "", fmt.Errorf("invalid session token")
	}

	if c.Purpose != purposeSession {
		return 0, "", fmt.Errorf("invalid token purpose")
	}

	// Enforce single-use: reject if this JTI was already consumed.
	expiry := c.ExpiresAt.Time
	if _, loaded := s.usedSessions.LoadOrStore(c.JTI, expiry); loaded {
		return 0, "", fmt.Errorf("session token already used")
	}

	// Lazily evict expired JTIs to prevent the map from growing unboundedly.
	s.evictExpiredSessions()

	return c.UserID, c.Next, nil
}

// evictExpiredSessions removes JTIs whose associated session has already expired.
// Called lazily on each ValidateSessionToken to avoid a background goroutine.
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

	// Lazily evict expired entries to bound memory usage.
	s.evictUsedTOTPCodes()

	return true, nil
}

// evictUsedTOTPCodes removes entries whose replay window has elapsed.
func (s *AuthService) evictUsedTOTPCodes() {
	now := time.Now()
	s.usedTOTPCodes.Range(func(key, value any) bool {
		if exp, ok := value.(time.Time); ok && now.After(exp) {
			s.usedTOTPCodes.Delete(key)
		}
		return true
	})
}

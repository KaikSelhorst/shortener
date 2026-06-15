package dto

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email")
	}
	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if len(r.Password) > 72 {
		return errors.New("password must be at most 72 characters")
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenResponse is kept for internal use by issueTokenPair.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// AuthState is the unified auth state machine response.
//
// next == "complete" → authentication done; access_token and refresh_token are present.
// next == "totp"     → second factor required; session is present.
//
// Future factors (e.g. "sms", "webauthn") follow the same pattern.
type AuthState struct {
	Next         string `json:"next"`
	Session      string `json:"session,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
}

// TOTPValidateRequest is the payload for the second login step when TOTP is active.
type TOTPValidateRequest struct {
	Session string `json:"session"`
	Code    string `json:"code"`
}

// TOTPSetupResponse returns the otpauth:// URI and the human-readable secret to the frontend.
type TOTPSetupResponse struct {
	URI    string `json:"uri"`
	Secret string `json:"secret"`
}

// TOTPConfirmRequest confirms the setup with a valid authenticator code.
type TOTPConfirmRequest struct {
	Code string `json:"code"`
}

// TOTPDisableRequest disables TOTP by requiring the current authenticator code.
type TOTPDisableRequest struct {
	Code string `json:"code"`
}

package dto

import (
	"errors"
	"time"
)

var ValidScopes = map[string]bool{
	"links:create":      true,
	"links:read":        true,
	"links:update":      true,
	"links:delete":      true,
	"projects:create":   true,
	"projects:read":     true,
	"projects:update":   true,
	"projects:delete":   true,
	"webhooks:read":     true,
	"webhooks:create":   true,
	"webhooks:delete":   true,
	"*":                 true,
}

type CreateAPIKeyRequest struct {
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes"`
	ProjectID *int64   `json:"project_id"`
}

func (r *CreateAPIKeyRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if len(r.Scopes) == 0 {
		return errors.New("at least one scope is required")
	}
	for _, s := range r.Scopes {
		if !ValidScopes[s] {
			return errors.New("invalid scope: " + s)
		}
	}
	return nil
}

type APIKeyResponse struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	ProjectID  *int64     `json:"project_id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	Scopes     []string   `json:"scopes"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// CreateAPIKeyResponse includes the raw token — returned only at creation time.
type CreateAPIKeyResponse struct {
	APIKeyResponse
	Token string `json:"token"`
}

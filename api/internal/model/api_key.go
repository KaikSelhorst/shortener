package model

import "time"

type APIKey struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	ProjectID  *int64     `json:"project_id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	KeyHash    string     `json:"-"`
	Scopes     []string   `json:"scopes"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (k *APIKey) HasScope(scope string) bool {
	for _, s := range k.Scopes {
		if s == "*" || s == scope {
			return true
		}
	}
	return false
}

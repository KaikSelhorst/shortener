package dto

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	customCodeRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	reservedCodes = map[string]bool{
		"health": true, "auth": true, "api-keys": true, "projects": true,
	}
)

func validateURL(raw string) error {
	if raw == "" {
		return errors.New("url is required")
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil || u.Host == "" {
		return errors.New("url is invalid")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("url must use http or https")
	}
	return nil
}

// LinkRequest is used for both creating and updating a link.
type LinkRequest struct {
	URL         string     `json:"url"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	OgImage     *string    `json:"og_image"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CustomCode  *string    `json:"custom_code"`
}

func (r *LinkRequest) Validate() error {
	if err := validateURL(r.URL); err != nil {
		return err
	}
	if r.ExpiresAt != nil && !r.ExpiresAt.After(time.Now().UTC()) {
		return errors.New("expires_at must be in the future")
	}
	if r.CustomCode != nil && *r.CustomCode != "" {
		code := *r.CustomCode
		if len(code) < 3 {
			return errors.New("custom code must be at least 3 characters")
		}
		if len(code) > 50 {
			return errors.New("custom code must be at most 50 characters")
		}
		if !customCodeRe.MatchString(code) {
			return errors.New("custom code may only contain letters, numbers, hyphens, and underscores")
		}
		if reservedCodes[strings.ToLower(code)] {
			return errors.New("this short code is reserved")
		}
	}
	return nil
}

type LinkResponse struct {
	ID          int64      `json:"id"`
	ProjectID   int64      `json:"project_id"`
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	OgImage     *string    `json:"og_image"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	ShortURL    string     `json:"short_url"`
	TotalClicks int64      `json:"total_clicks"`
}

type ListLinksResponse struct {
	Data       []LinkResponse `json:"data"`
	NextCursor *string        `json:"next_cursor"`
	PrevCursor *string        `json:"prev_cursor"`
	Limit      int            `json:"limit"`
}

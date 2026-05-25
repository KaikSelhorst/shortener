package dto

import (
	"errors"
	"net/url"
	"time"
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
}

func (r *LinkRequest) Validate() error {
	if err := validateURL(r.URL); err != nil {
		return err
	}
	if r.ExpiresAt != nil && !r.ExpiresAt.After(time.Now().UTC()) {
		return errors.New("expires_at must be in the future")
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
}

type ListLinksResponse struct {
	Data       []LinkResponse `json:"data"`
	NextCursor *string        `json:"next_cursor"`
	PrevCursor *string        `json:"prev_cursor"`
	Limit      int            `json:"limit"`
}

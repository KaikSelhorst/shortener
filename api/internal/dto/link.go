package dto

import (
	"errors"
	"time"
)

type CreateLinkRequest struct {
	URL         string  `json:"url"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	OgImage     *string `json:"og_image"`
}

func (r *CreateLinkRequest) Validate() error {
	if r.URL == "" {
		return errors.New("url is required")
	}
	return nil
}

type UpdateLinkRequest struct {
	URL         string  `json:"url"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	OgImage     *string `json:"og_image"`
}

func (r *UpdateLinkRequest) Validate() error {
	if r.URL == "" {
		return errors.New("url is required")
	}
	return nil
}

type LinkResponse struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	OgImage     *string   `json:"og_image"`
	CreatedAt   time.Time `json:"created_at"`
	ShortURL    string    `json:"short_url"`
}

type ListLinksResponse struct {
	Data       []LinkResponse `json:"data"`
	NextCursor *string        `json:"next_cursor"`
	PrevCursor *string        `json:"prev_cursor"`
	Limit      int            `json:"limit"`
}

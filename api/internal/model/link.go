package model

import "time"

type Link struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	OriginalURL string    `json:"original_url"`
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	OgImage     *string   `json:"og_image"`
	CreatedAt   time.Time `json:"created_at"`
}

package model

import "time"

type Link struct {
	ID          int64      `json:"id"`
	ProjectID   int64      `json:"project_id"`
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	OgImage     *string    `json:"og_image"`
	ExpiresAt   *time.Time `json:"expires_at"`
	MaxClicks   *int64     `json:"max_clicks"`
	CreatedAt   time.Time  `json:"created_at"`
	TotalClicks int64      `json:"total_clicks"`
}

func (l *Link) IsExpired() bool {
	return l.ExpiresAt != nil && l.ExpiresAt.Before(time.Now().UTC())
}

func (l *Link) IsClickLimitReached() bool {
	return l.MaxClicks != nil && l.TotalClicks >= *l.MaxClicks
}

package model

import "time"

type Click struct {
	ID        int64     `json:"id"`
	LinkID    int64     `json:"link_id"`
	UserAgent *string   `json:"user_agent"`
	IPAddress *string   `json:"ip_address"`
	Referer   *string   `json:"referer"`
	CreatedAt time.Time `json:"created_at"`
}

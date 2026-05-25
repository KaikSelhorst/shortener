package model

import "time"

type Click struct {
	ID        int64
	LinkID    int64
	UserAgent *string
	IPAddress *string
	Referer   *string
	CreatedAt time.Time
}

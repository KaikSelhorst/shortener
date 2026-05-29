package model

import "time"

type Click struct {
	ID             int64
	LinkID         int64
	UserAgent      *string
	IPHash         *string // HMAC-SHA256(ip, IP_HASH_SECRET) — nil when IP is unavailable
	Referer        *string
	DeviceType     string // "mobile" | "desktop" | "tablet" | "bot" | "unknown"
	ReferrerSource string // "instagram" | "facebook" | "twitter" | "tiktok" | "linkedin" | "whatsapp" | "youtube" | "google" | "direct" | "other"
	CreatedAt      time.Time
}

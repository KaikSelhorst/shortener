package model

import "time"

type ClicksOverTime struct {
	Date  time.Time `json:"date"`
	Count int64     `json:"count"`
}

type DeviceBreakdown struct {
	Mobile  int64 `json:"mobile"`
	Desktop int64 `json:"desktop"`
	Tablet  int64 `json:"tablet"`
	Bot     int64 `json:"bot"`
	Unknown int64 `json:"unknown"`
}

type ReferrerBreakdown struct {
	Direct    int64 `json:"direct"`
	Instagram int64 `json:"instagram"`
	Facebook  int64 `json:"facebook"`
	Twitter   int64 `json:"twitter"`
	TikTok    int64 `json:"tiktok"`
	LinkedIn  int64 `json:"linkedin"`
	WhatsApp  int64 `json:"whatsapp"`
	YouTube   int64 `json:"youtube"`
	Google    int64 `json:"google"`
	Other     int64 `json:"other"`
}

type LinkAnalytics struct {
	LinkID       int64             `json:"link_id"`
	ShortCode    string            `json:"short_code"`
	TotalClicks  int64             `json:"total_clicks"`
	UniqueClicks int64             `json:"unique_clicks"`
	OverTime     []ClicksOverTime  `json:"over_time"`
	Devices      DeviceBreakdown   `json:"devices"`
	Referrers    ReferrerBreakdown `json:"referrers"`
}

type TopLink struct {
	ShortCode   string  `json:"short_code"`
	OriginalURL string  `json:"original_url"`
	Title       *string `json:"title"`
	TotalClicks int64   `json:"total_clicks"`
}

type ProjectAnalytics struct {
	TotalClicks  int64             `json:"total_clicks"`
	UniqueClicks int64             `json:"unique_clicks"`
	OverTime     []ClicksOverTime  `json:"over_time"`
	Devices      DeviceBreakdown   `json:"devices"`
	Referrers    ReferrerBreakdown `json:"referrers"`
	TopLinks     []TopLink         `json:"top_links"`
}

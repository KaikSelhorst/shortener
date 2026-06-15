package model

import (
	"encoding/json"
	"time"
)

type Webhook struct {
	ID        string    `json:"id"`
	ProjectID int64     `json:"project_id"`
	URL       string    `json:"url"`
	Secret    string    `json:"-"`
	Events    []string  `json:"events"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

type WebhookDelivery struct {
	ID             string          `json:"id"`
	WebhookID      string          `json:"webhook_id"`
	Event          string          `json:"event"`
	Payload        json.RawMessage `json:"payload"`
	Status         string          `json:"status"`
	Attempts       int             `json:"attempts"`
	ResponseStatus *int            `json:"response_status"`
	NextRetryAt    time.Time       `json:"next_retry_at"`
	CreatedAt      time.Time       `json:"created_at"`
}

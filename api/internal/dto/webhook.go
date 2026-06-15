package dto

import (
	"errors"
	"net/url"
	"time"
)

var ValidWebhookEvents = map[string]bool{
	"link.clicked": true,
	"link.created": true,
	"link.updated": true,
	"link.deleted": true,
}

type CreateWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

func (r *CreateWebhookRequest) Validate() error {
	if r.URL == "" {
		return errors.New("url is required")
	}
	parsed, err := url.ParseRequestURI(r.URL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return errors.New("url must be a valid http or https URL")
	}
	if len(r.Events) == 0 {
		return errors.New("at least one event is required")
	}
	for _, e := range r.Events {
		if !ValidWebhookEvents[e] {
			return errors.New("invalid event: " + e)
		}
	}
	return nil
}

type WebhookResponse struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateWebhookResponse includes the plaintext secret — returned only at creation time.
type CreateWebhookResponse struct {
	WebhookResponse
	Secret string `json:"secret"`
}

type WebhookDeliveryResponse struct {
	ID             string    `json:"id"`
	WebhookID      int64     `json:"webhook_id"`
	Event          string    `json:"event"`
	Status         string    `json:"status"`
	Attempts       int       `json:"attempts"`
	ResponseStatus *int      `json:"response_status"`
	NextRetryAt    time.Time `json:"next_retry_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type ListDeliveriesResponse struct {
	Data    []WebhookDeliveryResponse `json:"data"`
	HasMore bool                      `json:"has_more"`
	Page    int                       `json:"page"`
}

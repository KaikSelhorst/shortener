package dto

import "errors"

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

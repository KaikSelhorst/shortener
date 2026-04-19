package dto

import (
	"errors"
	"regexp"
)

var validName = regexp.MustCompile(`^[\p{L}0-9 _-]+$`)

type CreateProjectRequest struct {
	Name string `json:"name"`
}

func (r *CreateProjectRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}

	if !validName.MatchString(r.Name) {
		return errors.New("name can only contain letters, digits, spaces, hyphens and underscores")
	}

	return nil
}

type UpdateProjectRequest struct {
	Name string `json:"name"`
}

func (r *UpdateProjectRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}

	if !validName.MatchString(r.Name) {
		return errors.New("name can only contain letters, digits, spaces, hyphens and underscores")
	}

	return nil
}

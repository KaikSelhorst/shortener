package handler

import (
	"net/http"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (c *HealthHandler) Ok(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

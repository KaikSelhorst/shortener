package handler

import (
	"net/http"

	"go.uber.org/zap"
)

type RedirectHandler struct {
	logger *zap.SugaredLogger
}

func NewRedirectHandler(logger *zap.SugaredLogger) *RedirectHandler {
	return &RedirectHandler{logger: logger}
}

func (c *RedirectHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

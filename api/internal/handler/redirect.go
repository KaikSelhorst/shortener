package handler

import (
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

type RedirectHandler struct {
	linkRepository *repository.LinkRepository
}

func NewRedirectHandler(linkRepository *repository.LinkRepository) *RedirectHandler {
	return &RedirectHandler{linkRepository: linkRepository}
}

func (h *RedirectHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	link, err := h.linkRepository.GetByCode(r.Context(), code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

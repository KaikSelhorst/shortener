package handler

import (
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/repository"
)

type RedirectHandler struct {
	linkRepository *repository.LinkRepository
}

func NewRedirectHandler(linkRepository *repository.LinkRepository) *RedirectHandler {
	return &RedirectHandler{linkRepository: linkRepository}
}

func (h *RedirectHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	link, err := h.linkRepository.GetByID(r.Context(), 1)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

package handler

import (
	"net"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type RedirectHandler struct {
	linkRepository *repository.LinkRepository
	tracker        *service.TrackerService
}

func NewRedirectHandler(linkRepository *repository.LinkRepository, tracker *service.TrackerService) *RedirectHandler {
	return &RedirectHandler{linkRepository: linkRepository, tracker: tracker}
}

func (h *RedirectHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	link, err := h.linkRepository.GetByCode(r.Context(), code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)

	click := model.Click{LinkID: link.ID}

	if ua := r.UserAgent(); ua != "" {
		click.UserAgent = &ua
	}
	if ref := r.Referer(); ref != "" {
		click.Referer = &ref
	}
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil && ip != "" {
		click.IPAddress = &ip
	}

	h.tracker.Record(click)
}

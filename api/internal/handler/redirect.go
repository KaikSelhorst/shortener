package handler

import (
	"net"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/cache"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type RedirectHandler struct {
	linkRepository repository.LinkRepo
	tracker        *service.TrackerService
	cache          *cache.LinkCache
	ipHashSecret   string
}

func NewRedirectHandler(linkRepository repository.LinkRepo, tracker *service.TrackerService, cache *cache.LinkCache, ipHashSecret string) *RedirectHandler {
	return &RedirectHandler{
		linkRepository: linkRepository,
		tracker:        tracker,
		cache:          cache,
		ipHashSecret:   ipHashSecret,
	}
}

func (h *RedirectHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	link, ok := h.cache.Get(code)
	if !ok {
		var err error
		link, err = h.linkRepository.GetByCode(r.Context(), code)
		if err != nil {
			writeError(w, http.StatusNotFound, "link not found")
			return
		}
		// Only cache non-expired links; caching an expired link wastes memory
		// since it will always result in a 410 on every subsequent request.
		if !link.IsExpired() {
			h.cache.Set(link)
		}
	}

	if link.IsExpired() {
		writeError(w, http.StatusGone, "link has expired")
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)

	ua := r.UserAgent()
	ref := r.Referer()

	click := model.Click{
		LinkID:         link.ID,
		DeviceType:     service.ParseDeviceType(ua),
		ReferrerSource: service.ParseReferrerSource(ref),
	}

	if ua != "" {
		click.UserAgent = &ua
	}
	if ref != "" {
		click.Referer = &ref
	}
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil && ip != "" {
		hash := service.HashIP(ip, h.ipHashSecret)
		click.IPHash = &hash
	}

	h.tracker.Record(click)
}

package handler

import (
	"net"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/cache"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/KaikSelhorst/shortener/internal/sse"
)

type clickedPayload struct {
	Event     string `json:"event"`
	ProjectID int64  `json:"project_id"`
	LinkID    int64  `json:"link_id"`
	ShortCode string `json:"short_code"`
}

type RedirectHandler struct {
	linkRepository repository.LinkRepo
	tracker        *service.TrackerService
	webhookService *service.WebhookService
	linkCache      *cache.LinkCache
	analyticsCache *cache.AnalyticsCache
	hub            *sse.Hub
	ipHashSecret   string
}

func NewRedirectHandler(
	linkRepository repository.LinkRepo,
	tracker *service.TrackerService,
	webhookService *service.WebhookService,
	linkCache *cache.LinkCache,
	analyticsCache *cache.AnalyticsCache,
	hub *sse.Hub,
	ipHashSecret string,
) *RedirectHandler {
	return &RedirectHandler{
		linkRepository: linkRepository,
		tracker:        tracker,
		webhookService: webhookService,
		linkCache:      linkCache,
		analyticsCache: analyticsCache,
		hub:            hub,
		ipHashSecret:   ipHashSecret,
	}
}

func (h *RedirectHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	link, ok := h.linkCache.Get(code)
	if !ok {
		// Use GetByCodeWithStats so TotalClicks is available for the click-limit check.
		// Only links without max_clicks are cached; click-limited links always go through
		// this path so the count is fresh on every request.
		var err error
		link, err = h.linkRepository.GetByCodeWithStats(r.Context(), code)
		if err != nil {
			writeError(w, http.StatusNotFound, "link not found")
			return
		}
		if !link.IsExpired() && link.MaxClicks == nil {
			h.linkCache.Set(link)
		}
	}

	if link.IsExpired() {
		writeError(w, http.StatusGone, "link has expired")
		return
	}

	if link.IsClickLimitReached() {
		writeError(w, http.StatusGone, "link has reached its click limit")
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)

	ua := r.UserAgent()
	ref := r.Referer()

	click := model.Click{
		LinkID:         link.ID,
		DeviceType:     service.ParseDeviceType(ua),
		ReferrerSource: service.ParseReferrerSource(ref),
		Browser:        service.ParseBrowserName(ua),
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
	h.hub.Notify(link.ProjectID, link.ID, link.ShortCode)
	h.analyticsCache.InvalidateProject(link.ProjectID)
	h.analyticsCache.InvalidateLink(link.ID)

	h.webhookService.DispatchAsync(link.ProjectID, "link.clicked", clickedPayload{
		Event:     "link.clicked",
		ProjectID: link.ProjectID,
		LinkID:    link.ID,
		ShortCode: link.ShortCode,
	})
}

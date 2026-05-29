package handler

import (
	"net/http"
	"time"

	"github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

type AnalyticsHandler struct {
	analytics repository.AnalyticsRepo
	projects  repository.ProjectRepo
	links     repository.LinkRepo
}

func NewAnalyticsHandler(
	analytics repository.AnalyticsRepo,
	projects repository.ProjectRepo,
	links repository.LinkRepo,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		analytics: analytics,
		projects:  projects,
		links:     links,
	}
}

func (h *AnalyticsHandler) GetProjectAnalytics(w http.ResponseWriter, r *http.Request) {
	if _, isAPIKey := middleware.APIKeyFromContext(r.Context()); isAPIKey {
		writeError(w, http.StatusForbidden, "API keys cannot access analytics")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")
	project, err := h.projects.FindBySlug(r.Context(), slug)
	if err != nil {
		repoError(w, err, "project not found")
		return
	}
	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	since, until := parsePeriod(r)

	data, err := h.analytics.GetProjectAnalytics(r.Context(), project.ID, since, until)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch analytics")
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *AnalyticsHandler) GetLinkAnalytics(w http.ResponseWriter, r *http.Request) {
	if _, isAPIKey := middleware.APIKeyFromContext(r.Context()); isAPIKey {
		writeError(w, http.StatusForbidden, "API keys cannot access analytics")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")
	project, err := h.projects.FindBySlug(r.Context(), slug)
	if err != nil {
		repoError(w, err, "project not found")
		return
	}
	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	code := chi.URLParam(r, "code")
	link, err := h.links.GetByCode(r.Context(), code)
	if err != nil {
		repoError(w, err, "link not found")
		return
	}
	if link.ProjectID != project.ID {
		writeError(w, http.StatusNotFound, "link not found")
		return
	}

	since, until := parsePeriod(r)

	data, err := h.analytics.GetLinkAnalytics(r.Context(), link.ID, since, until)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch analytics")
		return
	}

	data.ShortCode = link.ShortCode
	writeJSON(w, http.StatusOK, data)
}

// parsePeriod reads the ?period= query param and returns the since/until window.
// Accepted values: "7d", "30d" (default), "90d".
func parsePeriod(r *http.Request) (since, until time.Time) {
	until = time.Now().UTC()
	switch r.URL.Query().Get("period") {
	case "7d":
		since = until.AddDate(0, 0, -7)
	case "90d":
		since = until.AddDate(0, 0, -90)
	default:
		since = until.AddDate(0, 0, -30)
	}
	return
}

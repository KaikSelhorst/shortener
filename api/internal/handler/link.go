package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/KaikSelhorst/shortener/internal/cache"
	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type linkClickedPayload struct {
	Event     string `json:"event"`
	ProjectID int64  `json:"project_id"`
	LinkID    int64  `json:"link_id"`
	ShortCode string `json:"short_code"`
}

type linkEventPayload struct {
	Event     string        `json:"event"`
	ProjectID int64         `json:"project_id"`
	Link      dto.LinkResponse `json:"link"`
}

type LinkHandler struct {
	linkRepository    repository.LinkRepo
	projectRepository repository.ProjectRepo
	shortcodeService  *service.ShortcodeService
	webhookService    *service.WebhookService
	cache             *cache.LinkCache
	baseURL           string
	cursorSecret      string
}

func NewLinkHandler(
	linkRepository repository.LinkRepo,
	projectRepository repository.ProjectRepo,
	shortcodeService *service.ShortcodeService,
	webhookService *service.WebhookService,
	cache *cache.LinkCache,
	baseURL, cursorSecret string,
) *LinkHandler {
	return &LinkHandler{
		linkRepository:    linkRepository,
		projectRepository: projectRepository,
		shortcodeService:  shortcodeService,
		webhookService:    webhookService,
		cache:             cache,
		baseURL:           baseURL,
		cursorSecret:      cursorSecret,
	}
}

func (h *LinkHandler) isSelfReferential(rawURL string) bool {
	base, err := url.Parse(h.baseURL)
	if err != nil {
		return false
	}
	target, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return strings.EqualFold(target.Host, base.Host)
}

func (h *LinkHandler) toLinkResponse(link *model.Link) dto.LinkResponse {
	return dto.LinkResponse{
		ID:          link.ID,
		ProjectID:   link.ProjectID,
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		Title:       link.Title,
		Description: link.Description,
		OgImage:     link.OgImage,
		ExpiresAt:   link.ExpiresAt,
		CreatedAt:   link.CreatedAt,
		ShortURL:    h.baseURL + "/" + link.ShortCode,
		TotalClicks: link.TotalClicks,
	}
}

// assertProjectAccess checks ownership and API key scope for the given project.
// It writes an appropriate HTTP error and returns false if access is denied.
func (h *LinkHandler) assertProjectAccess(w http.ResponseWriter, r *http.Request, project *model.Project, userID int64) bool {
	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return false
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		writeError(w, http.StatusForbidden, "key not authorized for this project")
		return false
	}
	return true
}

// loadProjectBySlug resolves the {slug} URL parameter to a project and
// verifies ownership and API key scope. Returns the project and true on
// success, or writes an HTTP error and returns nil, false on failure.
func (h *LinkHandler) loadProjectBySlug(w http.ResponseWriter, r *http.Request, userID int64) (*model.Project, bool) {
	slug := chi.URLParam(r, "slug")
	project, err := h.projectRepository.FindBySlug(r.Context(), slug)
	if err != nil {
		repoError(w, err, "project not found")
		return nil, false
	}
	if !h.assertProjectAccess(w, r, project, userID) {
		return nil, false
	}
	return project, true
}

// loadProjectByID resolves a project by its numeric ID and verifies ownership
// and API key scope. Returns the project and true on success, or writes an
// HTTP error and returns nil, false on failure.
func (h *LinkHandler) loadProjectByID(w http.ResponseWriter, r *http.Request, userID int64, projectID int64) (*model.Project, bool) {
	project, err := h.projectRepository.GetByID(r.Context(), projectID)
	if err != nil {
		repoError(w, err, "project not found")
		return nil, false
	}
	if !h.assertProjectAccess(w, r, project, userID) {
		return nil, false
	}
	return project, true
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.LinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if h.isSelfReferential(req.URL) {
		writeError(w, http.StatusBadRequest, "url cannot point to this shortener")
		return
	}

	project, ok := h.loadProjectBySlug(w, r, userID)
	if !ok {
		return
	}

	newLink := &model.Link{
		ProjectID:   project.ID,
		OriginalURL: req.URL,
		Title:       req.Title,
		Description: req.Description,
		OgImage:     req.OgImage,
		ExpiresAt:   req.ExpiresAt,
	}

	if err := h.linkRepository.Create(r.Context(), newLink, h.shortcodeService.GenerateShortCode); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create link")
		return
	}

	resp := h.toLinkResponse(newLink)
	writeJSON(w, http.StatusCreated, resp)

	h.webhookService.DispatchAsync(project.ID, "link.created", linkEventPayload{
		Event:     "link.created",
		ProjectID: project.ID,
		Link:      resp,
	})
}

func (h *LinkHandler) ListLinks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	project, ok := h.loadProjectBySlug(w, r, userID)
	if !ok {
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var cursorCode string
	var direction string

	if raw := r.URL.Query().Get("cursor"); raw != "" {
		var err error
		direction, cursorCode, err = service.DecodeCursor(raw, h.cursorSecret)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid cursor")
			return
		}
	}

	cursorID := h.shortcodeService.DecodeShortCode(cursorCode)

	links, err := h.linkRepository.List(r.Context(), project.ID, cursorID, direction, limit+1)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list links")
		return
	}

	hasMore := len(links) > limit
	if hasMore {
		links = links[:limit]
	}

	var nextCursor, prevCursor *string

	if len(links) > 0 {
		shouldSetNext := hasMore || direction == "prev"
		shouldSetPrev := direction == "next" || (direction == "prev" && hasMore)

		if shouldSetNext {
			enc, err := service.EncodeCursor("next", links[len(links)-1].ShortCode, h.cursorSecret)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to encode cursor")
				return
			}
			nextCursor = &enc
		}

		if shouldSetPrev {
			enc, err := service.EncodeCursor("prev", links[0].ShortCode, h.cursorSecret)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to encode cursor")
				return
			}
			prevCursor = &enc
		}
	}

	data := make([]dto.LinkResponse, len(links))
	for i, link := range links {
		data[i] = h.toLinkResponse(link)
	}

	writeJSON(w, http.StatusOK, dto.ListLinksResponse{
		Data:       data,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		Limit:      limit,
	})
}

func (h *LinkHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	code := chi.URLParam(r, "code")

	link, err := h.linkRepository.GetByCodeWithStats(r.Context(), code)
	if err != nil {
		repoError(w, err, "link not found")
		return
	}

	if _, ok := h.loadProjectByID(w, r, userID, link.ProjectID); !ok {
		return
	}

	writeJSON(w, http.StatusOK, h.toLinkResponse(link))
}

func (h *LinkHandler) UpdateLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	code := chi.URLParam(r, "code")

	var req dto.LinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if h.isSelfReferential(req.URL) {
		writeError(w, http.StatusBadRequest, "url cannot point to this shortener")
		return
	}

	link, err := h.linkRepository.GetByCode(r.Context(), code)
	if err != nil {
		repoError(w, err, "link not found")
		return
	}

	if _, ok := h.loadProjectByID(w, r, userID, link.ProjectID); !ok {
		return
	}

	link.OriginalURL = req.URL
	link.Title = req.Title
	link.Description = req.Description
	link.OgImage = req.OgImage
	link.ExpiresAt = req.ExpiresAt

	if err := h.linkRepository.Update(r.Context(), link); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update link")
		return
	}

	h.cache.Delete(code)

	resp := h.toLinkResponse(link)
	writeJSON(w, http.StatusOK, resp)

	h.webhookService.DispatchAsync(link.ProjectID, "link.updated", linkEventPayload{
		Event:     "link.updated",
		ProjectID: link.ProjectID,
		Link:      resp,
	})
}

func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	code := chi.URLParam(r, "code")

	link, err := h.linkRepository.GetByCode(r.Context(), code)
	if err != nil {
		repoError(w, err, "link not found")
		return
	}

	if _, ok := h.loadProjectByID(w, r, userID, link.ProjectID); !ok {
		return
	}

	if err := h.linkRepository.DeleteByCode(r.Context(), code); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete link")
		return
	}

	h.cache.Delete(code)
	w.WriteHeader(http.StatusNoContent)

	h.webhookService.DispatchAsync(link.ProjectID, "link.deleted", linkEventPayload{
		Event:     "link.deleted",
		ProjectID: link.ProjectID,
		Link:      h.toLinkResponse(link),
	})
}

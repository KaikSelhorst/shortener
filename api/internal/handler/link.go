package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KaikSelhorst/shortener/internal/cache"
	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	linkRepository    *repository.LinkRepository
	projectRepository *repository.ProjectRepository
	cache             *cache.LinkCache
	baseURL           string
	cursorSecret      string
}

func NewLinkHandler(linkRepository *repository.LinkRepository, projectRepository *repository.ProjectRepository, cache *cache.LinkCache, baseURL, cursorSecret string) *LinkHandler {
	return &LinkHandler{
		linkRepository:    linkRepository,
		projectRepository: projectRepository,
		cache:             cache,
		baseURL:           baseURL,
		cursorSecret:      cursorSecret,
	}
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
	}
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	slug := chi.URLParam(r, "slug")
	project, err := h.projectRepository.FindBySlug(r.Context(), slug)
	if err != nil {
		repoError(w, err, "project not found")
		return
	}

	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		writeError(w, http.StatusForbidden, "key not authorized for this project")
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

	if err := h.linkRepository.Create(r.Context(), newLink, service.GenerateShortCode); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create link")
		return
	}

	writeJSON(w, http.StatusCreated, h.toLinkResponse(newLink))
}

func (h *LinkHandler) ListLinks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")
	project, err := h.projectRepository.FindBySlug(r.Context(), slug)
	if err != nil {
		repoError(w, err, "project not found")
		return
	}

	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		writeError(w, http.StatusForbidden, "key not authorized for this project")
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

	cursorID := service.DecodeShortCode(cursorCode)

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

	link, err := h.linkRepository.GetByCode(r.Context(), code)
	if err != nil {
		repoError(w, err, "link not found")
		return
	}

	project, err := h.projectRepository.GetByID(r.Context(), link.ProjectID)
	if err != nil {
		repoError(w, err, "project not found")
		return
	}
	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		writeError(w, http.StatusForbidden, "key not authorized for this project")
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

	var req dto.UpdateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	link, err := h.linkRepository.GetByCode(r.Context(), code)
	if err != nil {
		repoError(w, err, "link not found")
		return
	}

	project, err := h.projectRepository.GetByID(r.Context(), link.ProjectID)
	if err != nil {
		repoError(w, err, "project not found")
		return
	}
	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		writeError(w, http.StatusForbidden, "key not authorized for this project")
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

	writeJSON(w, http.StatusOK, h.toLinkResponse(link))
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

	project, err := h.projectRepository.GetByID(r.Context(), link.ProjectID)
	if err != nil {
		repoError(w, err, "project not found")
		return
	}
	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		writeError(w, http.StatusForbidden, "key not authorized for this project")
		return
	}

	if err := h.linkRepository.DeleteByCode(r.Context(), code); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete link")
		return
	}

	h.cache.Delete(code)

	w.WriteHeader(http.StatusNoContent)
}

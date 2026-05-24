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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slug := chi.URLParam(r, "slug")
	project, err := h.projectRepository.FindBySlug(r.Context(), slug)

	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	if project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		http.Error(w, "Forbidden: key not authorized for this project", http.StatusForbidden)
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
		http.Error(w, "failed to create link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.toLinkResponse(newLink))
}

func (h *LinkHandler) ListLinks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	slug := chi.URLParam(r, "slug")
	project, err := h.projectRepository.FindBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	if project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		http.Error(w, "Forbidden: key not authorized for this project", http.StatusForbidden)
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
			http.Error(w, "invalid cursor", http.StatusBadRequest)
			return
		}
	}

	cursorID := service.DecodeShortCode(cursorCode)

	links, err := h.linkRepository.List(r.Context(), project.ID, cursorID, direction, limit+1)
	if err != nil {
		http.Error(w, "failed to list links", http.StatusInternalServerError)
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
				http.Error(w, "failed to encode cursor", http.StatusInternalServerError)
				return
			}
			nextCursor = &enc
		}

		if shouldSetPrev {
			enc, err := service.EncodeCursor("prev", links[0].ShortCode, h.cursorSecret)
			if err != nil {
				http.Error(w, "failed to encode cursor", http.StatusInternalServerError)
				return
			}
			prevCursor = &enc
		}
	}

	data := make([]dto.LinkResponse, len(links))
	for i, link := range links {
		data[i] = h.toLinkResponse(link)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ListLinksResponse{
		Data:       data,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		Limit:      limit,
	})
}

func (h *LinkHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	code := chi.URLParam(r, "code")

	link, err := h.linkRepository.GetByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepository.GetByID(r.Context(), link.ProjectID)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}
	if project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		http.Error(w, "Forbidden: key not authorized for this project", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.toLinkResponse(link))
}

func (h *LinkHandler) UpdateLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	code := chi.URLParam(r, "code")

	var req dto.UpdateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link, err := h.linkRepository.GetByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepository.GetByID(r.Context(), link.ProjectID)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}
	if project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		http.Error(w, "Forbidden: key not authorized for this project", http.StatusForbidden)
		return
	}

	link.OriginalURL = req.URL
	link.Title = req.Title
	link.Description = req.Description
	link.OgImage = req.OgImage
	link.ExpiresAt = req.ExpiresAt

	if err := h.linkRepository.Update(r.Context(), link); err != nil {
		http.Error(w, "failed to update link", http.StatusInternalServerError)
		return
	}

	h.cache.Delete(code)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.toLinkResponse(link))
}

func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	code := chi.URLParam(r, "code")

	link, err := h.linkRepository.GetByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepository.GetByID(r.Context(), link.ProjectID)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}
	if project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if !middleware.ProjectAllowed(r.Context(), project.ID) {
		http.Error(w, "Forbidden: key not authorized for this project", http.StatusForbidden)
		return
	}

	if err := h.linkRepository.DeleteByCode(r.Context(), code); err != nil {
		http.Error(w, "failed to delete link", http.StatusInternalServerError)
		return
	}

	h.cache.Delete(code)

	w.WriteHeader(http.StatusNoContent)
}

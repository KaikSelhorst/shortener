package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	linkRepository    *repository.LinkRepository
	projectRepository *repository.ProjectRepository
}

func NewLinkHandler(linkRepository *repository.LinkRepository, projectRepository *repository.ProjectRepository) *LinkHandler {
	return &LinkHandler{
		linkRepository:    linkRepository,
		projectRepository: projectRepository,
	}
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
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

	newLink := &model.Link{
		ProjectID:   project.ID,
		OriginalURL: req.URL,
		Title:       req.Title,
		Description: req.Description,
		OgImage:     req.OgImage,
	}

	if err := h.linkRepository.Create(r.Context(), newLink); err != nil {
		http.Error(w, "failed to create link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newLink)
}

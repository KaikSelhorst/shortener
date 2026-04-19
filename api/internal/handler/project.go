package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	projectRepository *repository.ProjectRepository
}

func NewProjectHandler(projectRepository *repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{projectRepository: projectRepository}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project dto.CreateProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := project.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newProject := &model.Project{
		Name: project.Name,
		Slug: service.GenerateSlug(project.Name),
	}

	if err := h.projectRepository.Create(r.Context(), newProject); err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProject)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var project dto.UpdateProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := project.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slug := chi.URLParam(r, "slug")
	existingProject, err := h.projectRepository.FindBySlug(r.Context(), slug)

	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	existingProject.Name = project.Name
	existingProject.Slug = service.GenerateSlug(project.Name)

	if err := h.projectRepository.Update(r.Context(), existingProject); err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingProject)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	existingProject, err := h.projectRepository.FindBySlug(r.Context(), slug)

	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	if err := h.projectRepository.Delete(r.Context(), existingProject.ID); err != nil {
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

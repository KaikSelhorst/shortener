package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
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

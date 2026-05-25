package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgerrcode"
)

type ProjectHandler struct {
	projectRepository *repository.ProjectRepository
}

func NewProjectHandler(projectRepository *repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{projectRepository: projectRepository}
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	projects, err := h.projectRepository.FindAllByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list projects")
		return
	}

	if projects == nil {
		projects = []*model.Project{}
	}

	writeJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	newProject := &model.Project{
		UserID: userID,
		Name:   req.Name,
		Slug:   service.GenerateSlug(req.Name),
	}
	if err := h.projectRepository.Create(r.Context(), newProject); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			writeError(w, http.StatusConflict, "a project with this name already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create project")
		return
	}

	writeJSON(w, http.StatusCreated, newProject)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.UpdateProjectRequest
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

	project.Name = req.Name
	project.Slug = service.GenerateSlug(req.Name)

	if err := h.projectRepository.Update(r.Context(), project); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			writeError(w, http.StatusConflict, "a project with this name already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update project")
		return
	}

	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
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

	if err := h.projectRepository.Delete(r.Context(), project.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

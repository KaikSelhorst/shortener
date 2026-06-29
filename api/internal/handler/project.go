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
	"github.com/KaikSelhorst/shortener/internal/sse"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgerrcode"
)

type ProjectHandler struct {
	projectRepository repository.ProjectRepo
	hub               *sse.Hub
}

func NewProjectHandler(projectRepository repository.ProjectRepo, hub *sse.Hub) *ProjectHandler {
	return &ProjectHandler{projectRepository: projectRepository, hub: hub}
}

func toProjectResponse(p *model.Project) dto.ProjectResponse {
	return dto.ProjectResponse{
		ID:        p.ID,
		Name:      p.Name,
		Slug:      p.Slug,
		CreatedAt: p.CreatedAt,
	}
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

	data := make([]dto.ProjectResponse, len(projects))
	for i, p := range projects {
		data[i] = toProjectResponse(p)
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.ProjectRequest
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

	h.hub.InvalidateUserBootstrap(userID)
	writeJSON(w, http.StatusCreated, toProjectResponse(newProject))
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	slug := r.PathValue("slug")
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

	writeJSON(w, http.StatusOK, toProjectResponse(project))
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slug := r.PathValue("slug")
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

	h.hub.UnregisterProject(project.ID)
	w.WriteHeader(http.StatusNoContent)
}

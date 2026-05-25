package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type APIKeyHandler struct {
	apiKeyRepository *repository.APIKeyRepository
	authService      *service.AuthService
}

func NewAPIKeyHandler(apiKeyRepository *repository.APIKeyRepository, authService *service.AuthService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyRepository: apiKeyRepository,
		authService:      authService,
	}
}

func (h *APIKeyHandler) toResponse(k *model.APIKey) dto.APIKeyResponse {
	return dto.APIKeyResponse{
		ID:         k.ID,
		UserID:     k.UserID,
		ProjectID:  k.ProjectID,
		Name:       k.Name,
		KeyPrefix:  k.KeyPrefix,
		Scopes:     k.Scopes,
		LastUsedAt: k.LastUsedAt,
		CreatedAt:  k.CreatedAt,
	}
}

func (h *APIKeyHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	if _, isAPIKey := middleware.APIKeyFromContext(r.Context()); isAPIKey {
		writeError(w, http.StatusForbidden, "API keys cannot manage API keys")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw, hash, err := h.authService.GenerateAPIKey()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate key")
		return
	}

	key := &model.APIKey{
		UserID:    userID,
		ProjectID: req.ProjectID,
		Name:      req.Name,
		KeyPrefix: raw[:8],
		KeyHash:   hash,
		Scopes:    req.Scopes,
	}

	if err := h.apiKeyRepository.Create(r.Context(), key); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create api key")
		return
	}

	writeJSON(w, http.StatusCreated, dto.CreateAPIKeyResponse{
		APIKeyResponse: h.toResponse(key),
		Token:          raw,
	})
}

func (h *APIKeyHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	if _, isAPIKey := middleware.APIKeyFromContext(r.Context()); isAPIKey {
		writeError(w, http.StatusForbidden, "API keys cannot manage API keys")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	keys, err := h.apiKeyRepository.List(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list api keys")
		return
	}

	data := make([]dto.APIKeyResponse, len(keys))
	for i, k := range keys {
		data[i] = h.toResponse(k)
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *APIKeyHandler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	if _, isAPIKey := middleware.APIKeyFromContext(r.Context()); isAPIKey {
		writeError(w, http.StatusForbidden, "API keys cannot manage API keys")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.apiKeyRepository.Delete(r.Context(), id, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "api key not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete api key")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

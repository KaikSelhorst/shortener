package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type ApiKeyHandler struct {
	apiKeyRepository *repository.ApiKeyRepository
}

func NewApiKeyHandler(apiKeyRepository *repository.ApiKeyRepository) *ApiKeyHandler {
	return &ApiKeyHandler{apiKeyRepository: apiKeyRepository}
}

func (h *ApiKeyHandler) toResponse(k *model.ApiKey) dto.ApiKeyResponse {
	return dto.ApiKeyResponse{
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

func (h *ApiKeyHandler) CreateApiKey(w http.ResponseWriter, r *http.Request) {
	if _, isApiKey := middleware.ApiKeyFromContext(r.Context()); isApiKey {
		http.Error(w, "Forbidden: API keys cannot manage API keys", http.StatusForbidden)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	raw, hash, err := service.GenerateApiKey()
	if err != nil {
		http.Error(w, "failed to generate key", http.StatusInternalServerError)
		return
	}

	key := &model.ApiKey{
		UserID:    userID,
		ProjectID: req.ProjectID,
		Name:      req.Name,
		KeyPrefix: raw[:8],
		KeyHash:   hash,
		Scopes:    req.Scopes,
	}

	if err := h.apiKeyRepository.Create(r.Context(), key); err != nil {
		http.Error(w, "failed to create api key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.CreateApiKeyResponse{
		ApiKeyResponse: h.toResponse(key),
		Token:          raw,
	})
}

func (h *ApiKeyHandler) ListApiKeys(w http.ResponseWriter, r *http.Request) {
	if _, isApiKey := middleware.ApiKeyFromContext(r.Context()); isApiKey {
		http.Error(w, "Forbidden: API keys cannot manage API keys", http.StatusForbidden)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	keys, err := h.apiKeyRepository.List(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to list api keys", http.StatusInternalServerError)
		return
	}

	data := make([]dto.ApiKeyResponse, len(keys))
	for i, k := range keys {
		data[i] = h.toResponse(k)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *ApiKeyHandler) DeleteApiKey(w http.ResponseWriter, r *http.Request) {
	if _, isApiKey := middleware.ApiKeyFromContext(r.Context()); isApiKey {
		http.Error(w, "Forbidden: API keys cannot manage API keys", http.StatusForbidden)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.apiKeyRepository.Delete(r.Context(), id, userID); err != nil {
		if err == repository.ErrNotFound {
			http.Error(w, "api key not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete api key", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

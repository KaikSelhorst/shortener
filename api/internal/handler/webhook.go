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
	"github.com/KaikSelhorst/shortener/internal/uuidv7"
)

type WebhookHandler struct {
	webhookService    *service.WebhookService
	webhookRepository repository.WebhookRepo
	projectRepository repository.ProjectRepo
}

func NewWebhookHandler(
	webhookService *service.WebhookService,
	webhookRepository repository.WebhookRepo,
	projectRepository repository.ProjectRepo,
) *WebhookHandler {
	return &WebhookHandler{
		webhookService:    webhookService,
		webhookRepository: webhookRepository,
		projectRepository: projectRepository,
	}
}

func (h *WebhookHandler) toResponse(w *model.Webhook) dto.WebhookResponse {
	return dto.WebhookResponse{
		ID:        w.ID,
		ProjectID: w.ProjectID,
		URL:       w.URL,
		Events:    w.Events,
		Enabled:   w.Enabled,
		CreatedAt: w.CreatedAt,
	}
}

func (h *WebhookHandler) toDeliveryResponse(d *model.WebhookDelivery) dto.WebhookDeliveryResponse {
	return dto.WebhookDeliveryResponse{
		ID:             d.ID,
		WebhookID:      d.WebhookID,
		Event:          d.Event,
		Payload:        d.Payload,
		Status:         d.Status,
		Attempts:       d.Attempts,
		ResponseStatus: d.ResponseStatus,
		NextRetryAt:    d.NextRetryAt,
		CreatedAt:      d.CreatedAt,
	}
}

func (h *WebhookHandler) resolveProject(w http.ResponseWriter, r *http.Request, userID int64) (*model.Project, bool) {
	slug := r.PathValue("slug")
	project, err := h.projectRepository.FindBySlug(r.Context(), slug)
	if err != nil {
		repoError(w, err, "project not found")
		return nil, false
	}
	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return nil, false
	}
	return project, true
}

func (h *WebhookHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	project, ok := h.resolveProject(w, r, userID)
	if !ok {
		return
	}

	var req dto.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	plainSecret, err := service.GenerateWebhookSecret()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate secret")
		return
	}

	encryptedSecret, err := h.webhookService.EncryptSecret(plainSecret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to encrypt secret")
		return
	}

	wh := &model.Webhook{
		ProjectID: project.ID,
		URL:       req.URL,
		Secret:    encryptedSecret,
		Events:    req.Events,
		Enabled:   true,
	}

	if err := h.webhookRepository.Create(r.Context(), wh); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create webhook")
		return
	}

	writeJSON(w, http.StatusCreated, dto.CreateWebhookResponse{
		WebhookResponse: h.toResponse(wh),
		Secret:          plainSecret,
	})
}

func (h *WebhookHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	project, ok := h.resolveProject(w, r, userID)
	if !ok {
		return
	}

	webhooks, err := h.webhookRepository.List(r.Context(), project.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list webhooks")
		return
	}

	data := make([]dto.WebhookResponse, len(webhooks))
	for i, wh := range webhooks {
		data[i] = h.toResponse(wh)
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	project, ok := h.resolveProject(w, r, userID)
	if !ok {
		return
	}

	id := r.PathValue("id")
	if !uuidv7.IsValid(id) {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.webhookRepository.Delete(r.Context(), id, project.ID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "webhook not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete webhook")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WebhookHandler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	project, ok := h.resolveProject(w, r, userID)
	if !ok {
		return
	}

	webhookID := r.PathValue("id")
	if !uuidv7.IsValid(webhookID) {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if _, err := h.webhookRepository.GetByID(r.Context(), webhookID, project.ID); err != nil {
		repoError(w, err, "webhook not found")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	deliveries, err := h.webhookRepository.ListDeliveries(r.Context(), webhookID, limit+1, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list deliveries")
		return
	}

	hasMore := len(deliveries) > limit
	if hasMore {
		deliveries = deliveries[:limit]
	}

	data := make([]dto.WebhookDeliveryResponse, len(deliveries))
	for i, d := range deliveries {
		data[i] = h.toDeliveryResponse(d)
	}

	writeJSON(w, http.StatusOK, dto.ListDeliveriesResponse{
		Data:    data,
		HasMore: hasMore,
		Page:    page,
	})
}

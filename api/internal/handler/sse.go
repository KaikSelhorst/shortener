package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/sse"
	"github.com/go-chi/chi/v5"
)

type SSEHandler struct {
	hub      *sse.Hub
	projects repository.ProjectRepo
	links    repository.LinkRepo
}

func NewSSEHandler(hub *sse.Hub, projects repository.ProjectRepo, links repository.LinkRepo) *SSEHandler {
	return &SSEHandler{hub: hub, projects: projects, links: links}
}

func (h *SSEHandler) HandleProjectStream(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")
	project, err := h.projects.FindBySlug(r.Context(), slug)
	if err != nil {
		repoError(w, err, "project not found")
		return
	}
	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	h.hub.RegisterProjectUser(project.ID, userID)

	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Time{})

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	fmt.Fprintf(w, ": connected\n\n")
	_ = rc.Flush()

	ch := h.hub.SubscribeProject(project.ID)
	defer h.hub.UnsubscribeProject(project.ID, ch)

	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(evt)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			_ = rc.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": ping\n\n")
			_ = rc.Flush()
		}
	}
}

func (h *SSEHandler) HandleLinkStream(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slug := chi.URLParam(r, "slug")
	project, err := h.projects.FindBySlug(r.Context(), slug)
	if err != nil {
		repoError(w, err, "project not found")
		return
	}
	if project.UserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	h.hub.RegisterProjectUser(project.ID, userID)

	code := chi.URLParam(r, "code")
	link, err := h.links.GetByCode(r.Context(), code)
	if err != nil {
		repoError(w, err, "link not found")
		return
	}
	if link.ProjectID != project.ID {
		writeError(w, http.StatusNotFound, "link not found")
		return
	}

	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Time{})

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	fmt.Fprintf(w, ": connected\n\n")
	_ = rc.Flush()

	ch := h.hub.SubscribeLink(link.ID)
	defer h.hub.UnsubscribeLink(link.ID, ch)

	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(evt)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			_ = rc.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": ping\n\n")
			_ = rc.Flush()
		}
	}
}

func (h *SSEHandler) HandleUserStream(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Register all user's projects so Notify can fan out to this stream.
	// Skip the DB query when the hub already has a complete mapping for this user.
	if !h.hub.IsUserBootstrapped(userID) {
		projects, err := h.projects.FindAllByUserID(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load projects")
			return
		}
		for _, p := range projects {
			h.hub.RegisterProjectUser(p.ID, userID)
		}
		h.hub.MarkUserBootstrapped(userID)
	}

	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Time{})

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	fmt.Fprintf(w, ": connected\n\n")
	_ = rc.Flush()

	ch := h.hub.SubscribeUser(userID)
	defer h.hub.UnsubscribeUser(userID, ch)

	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(evt)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			_ = rc.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": ping\n\n")
			_ = rc.Flush()
		}
	}
}

package sse

import "sync"

// ClickEvent is broadcast to all subscribers when a redirect click occurs.
type ClickEvent struct {
	ShortCode string `json:"short_code"`
	ProjectID int64  `json:"project_id"`
	LinkID    int64  `json:"link_id"`
}

// Hub manages SSE subscribers keyed by project ID, link ID, or user ID.
type Hub struct {
	mu                sync.RWMutex
	byProject         map[int64]map[chan ClickEvent]struct{}
	byLink            map[int64]map[chan ClickEvent]struct{}
	byUser            map[int64]map[chan ClickEvent]struct{}
	projectUser       map[int64]int64    // projectID → userID reverse lookup
	bootstrappedUsers map[int64]struct{} // userIDs whose all projects are already registered
}

func NewHub() *Hub {
	return &Hub{
		byProject:         make(map[int64]map[chan ClickEvent]struct{}),
		byLink:            make(map[int64]map[chan ClickEvent]struct{}),
		byUser:            make(map[int64]map[chan ClickEvent]struct{}),
		projectUser:       make(map[int64]int64),
		bootstrappedUsers: make(map[int64]struct{}),
	}
}

func (h *Hub) SubscribeProject(projectID int64) chan ClickEvent {
	ch := make(chan ClickEvent, 16)
	h.mu.Lock()
	if h.byProject[projectID] == nil {
		h.byProject[projectID] = make(map[chan ClickEvent]struct{})
	}
	h.byProject[projectID][ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) UnsubscribeProject(projectID int64, ch chan ClickEvent) {
	h.mu.Lock()
	delete(h.byProject[projectID], ch)
	if len(h.byProject[projectID]) == 0 {
		delete(h.byProject, projectID)
	}
	h.mu.Unlock()
	close(ch)
}

func (h *Hub) SubscribeLink(linkID int64) chan ClickEvent {
	ch := make(chan ClickEvent, 16)
	h.mu.Lock()
	if h.byLink[linkID] == nil {
		h.byLink[linkID] = make(map[chan ClickEvent]struct{})
	}
	h.byLink[linkID][ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) UnsubscribeLink(linkID int64, ch chan ClickEvent) {
	h.mu.Lock()
	delete(h.byLink[linkID], ch)
	if len(h.byLink[linkID]) == 0 {
		delete(h.byLink, linkID)
	}
	h.mu.Unlock()
	close(ch)
}

func (h *Hub) SubscribeUser(userID int64) chan ClickEvent {
	ch := make(chan ClickEvent, 16)
	h.mu.Lock()
	if h.byUser[userID] == nil {
		h.byUser[userID] = make(map[chan ClickEvent]struct{})
	}
	h.byUser[userID][ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) UnsubscribeUser(userID int64, ch chan ClickEvent) {
	h.mu.Lock()
	delete(h.byUser[userID], ch)
	if len(h.byUser[userID]) == 0 {
		delete(h.byUser, userID)
	}
	h.mu.Unlock()
	close(ch)
}

// RegisterProjectUser records which user owns a project so Notify can fan out to user-level streams.
func (h *Hub) RegisterProjectUser(projectID, userID int64) {
	h.mu.Lock()
	h.projectUser[projectID] = userID
	h.mu.Unlock()
}

// UnregisterProject removes the project→user mapping and invalidates the user's bootstrap
// so the next dashboard open re-queries the full project list. Call when a project is deleted.
func (h *Hub) UnregisterProject(projectID int64) {
	h.mu.Lock()
	if userID, ok := h.projectUser[projectID]; ok {
		delete(h.bootstrappedUsers, userID)
		delete(h.projectUser, projectID)
	}
	h.mu.Unlock()
}

// IsUserBootstrapped reports whether all projects for userID are already registered in the hub.
func (h *Hub) IsUserBootstrapped(userID int64) bool {
	h.mu.RLock()
	_, ok := h.bootstrappedUsers[userID]
	h.mu.RUnlock()
	return ok
}

// MarkUserBootstrapped records that all projects for userID have been registered.
func (h *Hub) MarkUserBootstrapped(userID int64) {
	h.mu.Lock()
	h.bootstrappedUsers[userID] = struct{}{}
	h.mu.Unlock()
}

// InvalidateUserBootstrap clears the bootstrap flag so the next dashboard open re-queries.
// Call when a project is created for this user.
func (h *Hub) InvalidateUserBootstrap(userID int64) {
	h.mu.Lock()
	delete(h.bootstrappedUsers, userID)
	h.mu.Unlock()
}

// Notify broadcasts to all project-level, link-level, and user-level subscribers.
func (h *Hub) Notify(projectID, linkID int64, shortCode string) {
	evt := ClickEvent{ShortCode: shortCode, ProjectID: projectID, LinkID: linkID}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.byProject[projectID] {
		select {
		case ch <- evt:
		default:
		}
	}
	for ch := range h.byLink[linkID] {
		select {
		case ch <- evt:
		default:
		}
	}
	if userID, ok := h.projectUser[projectID]; ok {
		for ch := range h.byUser[userID] {
			select {
			case ch <- evt:
			default:
			}
		}
	}
}

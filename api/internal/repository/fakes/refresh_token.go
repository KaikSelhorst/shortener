package fakes

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
)

type RefreshTokenRepo struct {
	mu          sync.RWMutex
	tokens      map[string]*model.RefreshToken // keyed by token_hash
	nextID      atomic.Int64
	ReturnError error
}

func NewRefreshTokenRepo() *RefreshTokenRepo {
	return &RefreshTokenRepo{tokens: make(map[string]*model.RefreshToken)}
}

func (r *RefreshTokenRepo) Create(_ context.Context, rt *model.RefreshToken) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	rt.ID = r.nextID.Add(1)
	copy := *rt
	r.tokens[copy.TokenHash] = &copy
	return nil
}

func (r *RefreshTokenRepo) RevokeIfActive(_ context.Context, hash string) (*model.RefreshToken, error) {
	if r.ReturnError != nil {
		return nil, r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	rt, ok := r.tokens[hash]
	if !ok {
		return nil, repository.ErrNotFound
	}
	copy := *rt
	delete(r.tokens, hash)
	return &copy, nil
}

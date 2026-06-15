package fakes

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
)

type LinkRepo struct {
	mu          sync.RWMutex
	links       map[string]*model.Link // keyed by short_code
	nextID      atomic.Int64
	ReturnError error
}

func NewLinkRepo() *LinkRepo {
	return &LinkRepo{links: make(map[string]*model.Link)}
}

func (r *LinkRepo) Create(_ context.Context, link *model.Link, generateCode func(uint64) (string, error)) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	link.ID = r.nextID.Add(1)
	code, err := generateCode(uint64(link.ID))
	if err != nil {
		return err
	}
	link.ShortCode = code

	copy := *link
	r.links[copy.ShortCode] = &copy
	return nil
}

func (r *LinkRepo) GetByCode(_ context.Context, code string) (*model.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	l, ok := r.links[code]
	if !ok {
		return nil, repository.ErrNotFound
	}
	copy := *l
	return &copy, nil
}

func (r *LinkRepo) GetByCodeWithStats(_ context.Context, code string) (*model.Link, error) {
	return r.GetByCode(context.Background(), code)
}

func (r *LinkRepo) List(_ context.Context, projectID int64, _ uint64, _ string, limit int) ([]*model.Link, error) {
	if r.ReturnError != nil {
		return nil, r.ReturnError
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []*model.Link
	for _, l := range r.links {
		if l.ProjectID == projectID {
			copy := *l
			out = append(out, &copy)
			if len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

func (r *LinkRepo) Update(_ context.Context, link *model.Link) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.links[link.ShortCode]; !ok {
		return repository.ErrNotFound
	}
	copy := *link
	r.links[copy.ShortCode] = &copy
	return nil
}

func (r *LinkRepo) DeleteByCode(_ context.Context, code string) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.links, code)
	return nil
}

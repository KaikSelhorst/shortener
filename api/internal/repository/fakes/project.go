package fakes

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
)

type ProjectRepo struct {
	mu          sync.RWMutex
	projects    map[int64]*model.Project
	nextID      atomic.Int64
	ReturnError error
}

func NewProjectRepo() *ProjectRepo {
	return &ProjectRepo{projects: make(map[int64]*model.Project)}
}

func (r *ProjectRepo) Create(_ context.Context, p *model.Project) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	p.ID = r.nextID.Add(1)
	copy := *p
	r.projects[copy.ID] = &copy
	return nil
}

func (r *ProjectRepo) GetByID(_ context.Context, id int64) (*model.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.projects[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	copy := *p
	return &copy, nil
}

func (r *ProjectRepo) FindBySlug(_ context.Context, slug string) (*model.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.projects {
		if p.Slug == slug {
			copy := *p
			return &copy, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *ProjectRepo) FindAllByUserID(_ context.Context, userID int64) ([]*model.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []*model.Project
	for _, p := range r.projects {
		if p.UserID == userID {
			copy := *p
			out = append(out, &copy)
		}
	}
	return out, nil
}

func (r *ProjectRepo) Update(_ context.Context, p *model.Project) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.projects[p.ID]; !ok {
		return repository.ErrNotFound
	}
	copy := *p
	r.projects[copy.ID] = &copy
	return nil
}

func (r *ProjectRepo) Delete(_ context.Context, id int64) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.projects, id)
	return nil
}

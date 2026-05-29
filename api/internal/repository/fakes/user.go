package fakes

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	mu          sync.RWMutex
	users       map[int64]*model.User
	byEmail     map[string]int64
	nextID      atomic.Int64
	ReturnError error // returned by all mutating calls if set
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users:   make(map[int64]*model.User),
		byEmail: make(map[string]int64),
	}
}

func (r *UserRepo) Create(_ context.Context, user *model.User) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[user.Email]; exists {
		return &pgconn.PgError{Code: pgerrcode.UniqueViolation}
	}

	user.ID = r.nextID.Add(1)
	copy := *user
	r.users[copy.ID] = &copy
	r.byEmail[copy.Email] = copy.ID
	return nil
}

func (r *UserRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byEmail[email]
	if !ok {
		return nil, repository.ErrNotFound
	}
	copy := *r.users[id]
	return &copy, nil
}

func (r *UserRepo) FindByID(_ context.Context, id int64) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	copy := *u
	return &copy, nil
}

func (r *UserRepo) SaveTOTPSecret(_ context.Context, userID int64, secret string) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[userID]
	if !ok {
		return repository.ErrNotFound
	}
	u.TOTPSecret = &secret
	return nil
}

func (r *UserRepo) SetTOTPEnabled(_ context.Context, userID int64, enabled bool) error {
	if r.ReturnError != nil {
		return r.ReturnError
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[userID]
	if !ok {
		return repository.ErrNotFound
	}
	u.TOTPEnabled = enabled
	if !enabled {
		u.TOTPSecret = nil
	}
	return nil
}

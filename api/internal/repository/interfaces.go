package repository

import (
	"context"
	"time"

	"github.com/KaikSelhorst/shortener/internal/model"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	SaveTOTPSecret(ctx context.Context, userID int64, secret string) error
	SetTOTPEnabled(ctx context.Context, userID int64, enabled bool) error
}

type LinkRepo interface {
	Create(ctx context.Context, link *model.Link, generateCode func(uint64) (string, error)) error
	GetByCode(ctx context.Context, code string) (*model.Link, error)
	List(ctx context.Context, projectID int64, cursor uint64, direction string, limit int) ([]*model.Link, error)
	Update(ctx context.Context, link *model.Link) error
	DeleteByCode(ctx context.Context, code string) error
}

type ProjectRepo interface {
	Create(ctx context.Context, project *model.Project) error
	GetByID(ctx context.Context, id int64) (*model.Project, error)
	FindBySlug(ctx context.Context, slug string) (*model.Project, error)
	FindAllByUserID(ctx context.Context, userID int64) ([]*model.Project, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id int64) error
}

type RefreshTokenRepo interface {
	Create(ctx context.Context, rt *model.RefreshToken) error
	RevokeIfActive(ctx context.Context, hash string) (*model.RefreshToken, error)
}

type APIKeyRepo interface {
	Create(ctx context.Context, key *model.APIKey) error
	List(ctx context.Context, userID int64) ([]*model.APIKey, error)
	GetByHash(ctx context.Context, hash string) (*model.APIKey, error)
	Delete(ctx context.Context, id int64, userID int64) error
	UpdateLastUsed(ctx context.Context, id int64) error
}

type AnalyticsRepo interface {
	GetLinkAnalytics(ctx context.Context, linkID int64, since, until time.Time) (*model.LinkAnalytics, error)
	GetProjectAnalytics(ctx context.Context, projectID int64, since, until time.Time) (*model.ProjectAnalytics, error)
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.QueryRow(ctx,
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at",
		user.Email, user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx,
		"SELECT id, email, password_hash, totp_secret, totp_enabled, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.TOTPSecret, &user.TOTPEnabled, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx,
		"SELECT id, email, password_hash, totp_secret, totp_enabled, created_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.TOTPSecret, &user.TOTPEnabled, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// SaveTOTPSecret persists the pending TOTP secret generated during setup.
// totp_enabled remains false until the user confirms with a valid code.
func (r *UserRepository) SaveTOTPSecret(ctx context.Context, userID int64, secret string) error {
	tag, err := r.db.Exec(ctx,
		"UPDATE users SET totp_secret = $1 WHERE id = $2",
		secret, userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user %d not found", userID)
	}
	return nil
}

// SetTOTPEnabled activates or deactivates TOTP for the given user.
// When disabling, the stored secret is cleared as well.
func (r *UserRepository) SetTOTPEnabled(ctx context.Context, userID int64, enabled bool) error {
	var query string
	if enabled {
		query = "UPDATE users SET totp_enabled = TRUE WHERE id = $1"
	} else {
		query = "UPDATE users SET totp_enabled = FALSE, totp_secret = NULL WHERE id = $1"
	}
	tag, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user %d not found", userID)
	}
	return nil
}

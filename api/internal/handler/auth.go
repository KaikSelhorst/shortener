package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepository         *repository.UserRepository
	refreshTokenRepository *repository.RefreshTokenRepository
	authService            *service.AuthService
}

func NewAuthHandler(
	userRepository *repository.UserRepository,
	refreshTokenRepository *repository.RefreshTokenRepository,
	authService *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		authService:            authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: string(hash),
	}
	if err := h.userRepository.Create(r.Context(), user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userRepository.FindByEmail(r.Context(), req.Email)
	if err != nil {
		_, _ = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	tokens, err := h.issueTokenPair(r, user.ID)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.RefreshToken == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	hash := service.HashToken(req.RefreshToken)
	rt, err := h.refreshTokenRepository.RevokeIfActive(r.Context(), hash)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	tokens, err := h.issueTokenPair(r, rt.UserID)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.RefreshToken == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	hash := service.HashToken(req.RefreshToken)
	_, _ = h.refreshTokenRepository.RevokeIfActive(r.Context(), hash)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) issueTokenPair(r *http.Request, userID int64) (*dto.TokenResponse, error) {
	raw, hash, err := h.authService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	rt := &model.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(service.RefreshTokenTTL),
	}
	if err := h.refreshTokenRepository.Create(r.Context(), rt); err != nil {
		return nil, err
	}

	accessToken, err := h.authService.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: raw,
		TokenType:    "Bearer",
		ExpiresIn:    int(service.AccessTokenTTL.Seconds()),
	}, nil
}

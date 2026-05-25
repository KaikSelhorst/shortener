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
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: string(hash),
	}
	if err := h.userRepository.Create(r.Context(), user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			writeError(w, http.StatusConflict, "email already in use")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	// Issue tokens immediately so the client is logged in right after registration.
	tokens, err := h.issueTokenPair(r, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	writeJSON(w, http.StatusCreated, tokens)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userRepository.FindByEmail(r.Context(), req.Email)
	if err != nil {
		_, _ = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	tokens, err := h.issueTokenPair(r, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	hash := service.HashToken(req.RefreshToken)
	rt, err := h.refreshTokenRepository.RevokeIfActive(r.Context(), hash)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	tokens, err := h.issueTokenPair(r, rt.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "refresh_token is required")
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

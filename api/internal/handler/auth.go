package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepository         repository.UserRepo
	refreshTokenRepository repository.RefreshTokenRepo
	authService            *service.AuthService
	registrationDisabled   bool
}

func NewAuthHandler(
	userRepository repository.UserRepo,
	refreshTokenRepository repository.RefreshTokenRepo,
	authService *service.AuthService,
	registrationDisabled bool,
) *AuthHandler {
	return &AuthHandler{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		authService:            authService,
		registrationDisabled:   registrationDisabled,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if h.registrationDisabled {
		writeError(w, http.StatusForbidden, "registration is disabled")
		return
	}
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
	// New users never have TOTP enabled, so we go straight to "complete".
	tokens, err := h.issueTokenPair(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	writeJSON(w, http.StatusCreated, dto.AuthState{
		Next:         "complete",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
	})
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

	if user.TOTPEnabled {
		session, err := h.authService.GenerateSessionToken(user.ID, "totp")
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to generate session token")
			return
		}
		writeJSON(w, http.StatusOK, dto.AuthState{Next: "totp", Session: session})
		return
	}

	tokens, err := h.issueTokenPair(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	writeJSON(w, http.StatusOK, dto.AuthState{
		Next:         "complete",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
	})
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

	tokens, err := h.issueTokenPair(r.Context(), rt.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	writeJSON(w, http.StatusOK, dto.AuthState{
		Next:         "complete",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
	})
}

// Me handles GET /auth/me (requires auth).
// Returns the current user's public profile including totp_enabled status.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Block API keys: TOTP status is sensitive account information not needed
	// by programmatic API clients, and exposing it widens the attack surface.
	if _, isAPIKey := middleware.APIKeyFromContext(r.Context()); isAPIKey {
		writeError(w, http.StatusForbidden, "API keys cannot access account profile")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userRepository.FindByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":           user.ID,
		"email":        user.Email,
		"totp_enabled": user.TOTPEnabled,
		"created_at":   user.CreatedAt,
	})
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

	// Revoke the access token immediately so it cannot be reused within its
	// remaining TTL window. The token is optional — clients that omit it just
	// keep a short-lived window of access until natural expiry.
	if authHeader := r.Header.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		if !strings.HasPrefix(accessToken, "sk_") {
			_ = h.authService.RevokeAccessToken(accessToken)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) issueTokenPair(ctx context.Context, userID int64) (*dto.TokenResponse, error) {
	return issueTokenPairFor(ctx, h.refreshTokenRepository, h.authService, userID)
}

// issueTokenPairFor is a package-level helper shared by AuthHandler and TOTPHandler.
func issueTokenPairFor(
	ctx context.Context,
	rtRepo repository.RefreshTokenRepo,
	authSvc *service.AuthService,
	userID int64,
) (*dto.TokenResponse, error) {
	raw, hash, err := authSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	rt := &model.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(service.RefreshTokenTTL),
	}
	if err := rtRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	accessToken, err := authSvc.GenerateAccessToken(userID)
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

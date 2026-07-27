package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KaikSelhorst/shortener/internal/dto"
	"github.com/KaikSelhorst/shortener/internal/middleware"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
)

const totpIssuer = "Docut"

type TOTPHandler struct {
	userRepository repository.UserRepo
	authService    *service.AuthService
	// refreshTokenRepository is needed to issue the final token pair.
	refreshTokenRepository repository.RefreshTokenRepo
}

func NewTOTPHandler(
	userRepository repository.UserRepo,
	refreshTokenRepository repository.RefreshTokenRepo,
	authService *service.AuthService,
) *TOTPHandler {
	return &TOTPHandler{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		authService:            authService,
	}
}

// ValidateMFA handles POST /auth/mfa/totp.
// It is the second step of the login flow when TOTP is enabled.
// Expects { session, code }; issues real tokens on success.
//
// The session token is intentionally parsed without being consumed first so that
// a mistyped TOTP code does not burn the session — the user can retry within the
// 5-minute session window. The session is consumed only after TOTP succeeds.
func (h *TOTPHandler) ValidateMFA(w http.ResponseWriter, r *http.Request) {
	var req dto.TOTPValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if req.Session == "" || req.Code == "" {
		writeError(w, http.StatusBadRequest, "session and code are required")
		return
	}

	// Parse without consuming so a wrong TOTP code does not burn the session.
	userID, next, jti, expiry, err := h.authService.ParseSessionToken(req.Session)
	if err != nil || next != "totp" {
		writeError(w, http.StatusUnauthorized, "invalid or expired session")
		return
	}

	user, err := h.userRepository.FindByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid or expired session")
		return
	}

	if !user.TOTPEnabled || user.TOTPSecret == nil {
		writeError(w, http.StatusBadRequest, "totp is not enabled for this account")
		return
	}

	// ValidateTOTPAndConsume enforces single-use within the 90-second window,
	// preventing replay of an intercepted code during the same time step.
	ok, err := h.authService.ValidateTOTPAndConsume(userID, *user.TOTPSecret, req.Code)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "totp validation error")
		return
	}
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid totp code")
		return
	}

	// Consume the session only after successful TOTP validation.
	if err := h.authService.ConsumeSession(jti, expiry); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid or expired session")
		return
	}

	tokens, err := issueTokenPairFor(r.Context(), h.refreshTokenRepository, h.authService, userID)
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

// Setup handles POST /auth/totp/setup (requires auth).
// Generates a new TOTP secret, persists it (totp_enabled stays false),
// and returns the otpauth:// URI and the raw secret for manual entry.
func (h *TOTPHandler) Setup(w http.ResponseWriter, r *http.Request) {
	// Block API keys: only interactive (JWT) sessions may manage account security.
	if _, isAPIKey := middleware.APIKeyFromContext(r.Context()); isAPIKey {
		writeError(w, http.StatusForbidden, "API keys cannot manage account security settings")
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

	// Prevent overwriting an active secret — would lock the user out immediately.
	if user.TOTPEnabled {
		writeError(w, http.StatusConflict, "totp is already enabled — disable it before setting up again")
		return
	}

	secret, err := service.GenerateTOTPSecret()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate totp secret")
		return
	}

	if err := h.userRepository.SaveTOTPSecret(r.Context(), userID, secret); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save totp secret")
		return
	}

	uri := service.TOTPUri(secret, user.Email, totpIssuer)
	writeJSON(w, http.StatusOK, dto.TOTPSetupResponse{URI: uri, Secret: secret})
}

// Confirm handles POST /auth/totp/confirm (requires auth).
// Validates the code against the pending secret and activates TOTP.
func (h *TOTPHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	// Block API keys: only interactive (JWT) sessions may manage account security.
	if _, isAPIKey := middleware.APIKeyFromContext(r.Context()); isAPIKey {
		writeError(w, http.StatusForbidden, "API keys cannot manage account security settings")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.TOTPConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if req.Code == "" {
		writeError(w, http.StatusBadRequest, "code is required")
		return
	}

	user, err := h.userRepository.FindByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if user.TOTPSecret == nil {
		writeError(w, http.StatusBadRequest, "totp setup not initiated — call /auth/totp/setup first")
		return
	}

	if user.TOTPEnabled {
		writeError(w, http.StatusConflict, "totp is already enabled")
		return
	}

	valid, err2 := h.authService.ValidateTOTPAndConsume(userID, *user.TOTPSecret, req.Code)
	if err2 != nil {
		writeError(w, http.StatusInternalServerError, "totp validation error")
		return
	}
	if !valid {
		writeError(w, http.StatusUnauthorized, "invalid totp code")
		return
	}

	if err := h.userRepository.SetTOTPEnabled(r.Context(), userID, true); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to enable totp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Disable handles DELETE /auth/totp (requires auth).
// Validates the current TOTP code and deactivates TOTP, clearing the secret.
func (h *TOTPHandler) Disable(w http.ResponseWriter, r *http.Request) {
	// Block API keys: only interactive (JWT) sessions may manage account security.
	if _, isAPIKey := middleware.APIKeyFromContext(r.Context()); isAPIKey {
		writeError(w, http.StatusForbidden, "API keys cannot manage account security settings")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.TOTPDisableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if req.Code == "" {
		writeError(w, http.StatusBadRequest, "code is required")
		return
	}

	user, err := h.userRepository.FindByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if !user.TOTPEnabled || user.TOTPSecret == nil {
		writeError(w, http.StatusBadRequest, "totp is not enabled for this account")
		return
	}

	valid, err2 := h.authService.ValidateTOTPAndConsume(userID, *user.TOTPSecret, req.Code)
	if err2 != nil {
		writeError(w, http.StatusInternalServerError, "totp validation error")
		return
	}
	if !valid {
		writeError(w, http.StatusUnauthorized, "invalid totp code")
		return
	}

	if err := h.userRepository.SetTOTPEnabled(r.Context(), userID, false); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to disable totp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

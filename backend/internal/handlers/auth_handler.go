package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"backend/internal/logger"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo  repository.UserRepository
	logger    *logger.Logger
	jwtSecret string
}

func NewAuthHandler(userRepo repository.UserRepository, log *logger.Logger, jwtSecret string) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, logger: log, jwtSecret: jwtSecret}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	req.Username = trimSpace(req.Username)

	if len(req.Username) < 3 || len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "Username must be at least 3 chars, password at least 6 chars")
		return
	}
	if len(req.Password) > 72 {
		writeError(w, http.StatusBadRequest, "Password too long (max 72 chars)")
		return
	}

	existing, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		h.logger.WithError(err).Error("Register: find by username")
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if existing != nil {
		writeError(w, http.StatusConflict, "Username already taken")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.WithError(err).Error("Register: bcrypt error")
		writeError(w, http.StatusInternalServerError, "Server error")
		return
	}

	user, err := h.userRepo.Create(req.Username, string(hash), "user")
	if err != nil {
		h.logger.WithError(err).Error("Register: create user")
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		h.logger.WithError(err).Error("Register: token generation")
		writeError(w, http.StatusInternalServerError, "Server error")
		return
	}

	json.NewEncoder(w).Encode(models.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	req.Username = trimSpace(req.Username)

	user, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		h.logger.WithError(err).Error("Login: find by username")
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		h.logger.WithError(err).Error("Login: token generation")
		writeError(w, http.StatusInternalServerError, "Server error")
		return
	}

	h.logger.WithField("user_id", user.ID).Info("user logged in")

	json.NewEncoder(w).Encode(models.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	user, err := h.userRepo.FindByID(claims.UserID)
	if err != nil {
		h.logger.WithError(err).Error("Me: find by id")
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if user == nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) generateToken(user *models.User) (string, error) {
	claims := middleware.UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

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

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 || len(req.Password) < 6 {
		http.Error(w, `{"error":"Username must be at least 3 chars, password at least 6 chars"}`, http.StatusBadRequest)
		return
	}

	existing, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		h.logger.Errorf("Register: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	if existing != nil {
		http.Error(w, `{"error":"Username already taken"}`, http.StatusConflict)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Errorf("Register: bcrypt error: %v", err)
		http.Error(w, `{"error":"Server error"}`, http.StatusInternalServerError)
		return
	}

	user, err := h.userRepo.Create(req.Username, string(hash), "user")
	if err != nil {
		h.logger.Errorf("Register: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		h.logger.Errorf("Register: token error: %v", err)
		http.Error(w, `{"error":"Server error"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		h.logger.Errorf("Login: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, `{"error":"Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, `{"error":"Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		h.logger.Errorf("Login: token error: %v", err)
		http.Error(w, `{"error":"Server error"}`, http.StatusInternalServerError)
		return
	}

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
		http.Error(w, `{"error":"Not authenticated"}`, http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.FindByID(claims.UserID)
	if err != nil {
		h.logger.Errorf("Me: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
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

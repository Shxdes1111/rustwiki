package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"backend/internal/logger"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/repository"
)

var mimeToExt = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/webp": ".webp",
	"image/avif": ".avif",
}

func decodeDataURL(s string) ([]byte, string, error) {
	if !strings.HasPrefix(s, "data:") {
		return nil, "", fmt.Errorf("invalid data URL")
	}

	comma := strings.Index(s, ",")
	if comma < 0 {
		return nil, "", fmt.Errorf("invalid data URL")
	}

	header := s[5:comma]
	encoded := s[comma+1:]

	parts := strings.SplitN(header, ";", 2)
	if len(parts) != 2 || parts[1] != "base64" {
		return nil, "", fmt.Errorf("invalid data URL encoding")
	}
	mimeType := parts[0]

	ext, ok := mimeToExt[mimeType]
	if !ok {
		return nil, "", fmt.Errorf("unsupported image type: %s", mimeType)
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, "", err
	}

	if len(data) > maxUploadSize {
		return nil, "", fmt.Errorf("image too large (max %d bytes)", maxUploadSize)
	}

	return data, ext, nil
}

type SuggestionHandler struct {
	suggestionRepo repository.SuggestionRepository
	weaponRepo     repository.WeaponRepository
	logger         *logger.Logger
}

func NewSuggestionHandler(suggestionRepo repository.SuggestionRepository, weaponRepo repository.WeaponRepository, log *logger.Logger) *SuggestionHandler {
	return &SuggestionHandler{suggestionRepo: suggestionRepo, weaponRepo: weaponRepo, logger: log}
}

func (h *SuggestionHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		http.Error(w, `{"error":"Not authenticated"}`, http.StatusUnauthorized)
		return
	}

	var req models.CreateWeaponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if len(req.Description) > 500 {
		http.Error(w, `{"error":"Description too long (max 500 chars)"}`, http.StatusBadRequest)
		return
	}

	raw, err := json.Marshal(req)
	if err != nil {
		http.Error(w, `{"error":"Failed to encode payload"}`, http.StatusInternalServerError)
		return
	}
	payload := json.RawMessage(raw)

	s, err := h.suggestionRepo.Create(claims.UserID, payload)
	if err != nil {
		h.logger.Errorf("CreateSuggestion: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (h *SuggestionHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	suggestions, err := h.suggestionRepo.FindAll()
	if err != nil {
		h.logger.Errorf("ListSuggestions: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	if suggestions == nil {
		suggestions = []models.WeaponSuggestion{}
	}

	json.NewEncoder(w).Encode(suggestions)
}

func (h *SuggestionHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error":"Invalid suggestion ID"}`, http.StatusBadRequest)
		return
	}

	s, err := h.suggestionRepo.FindByID(id)
	if err != nil {
		h.logger.Errorf("GetSuggestion: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	if s == nil {
		http.Error(w, `{"error":"Suggestion not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(s)
}

func (h *SuggestionHandler) Approve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error":"Invalid suggestion ID"}`, http.StatusBadRequest)
		return
	}

	claims := middleware.GetUserClaims(r)

	s, err := h.suggestionRepo.FindByID(id)
	if err != nil {
		h.logger.Errorf("ApproveSuggestion: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	if s == nil {
		http.Error(w, `{"error":"Suggestion not found"}`, http.StatusNotFound)
		return
	}
	if s.Status != "pending" {
		http.Error(w, `{"error":"Suggestion already reviewed"}`, http.StatusBadRequest)
		return
	}

	var req models.CreateWeaponRequest
	if err := json.Unmarshal(s.Payload, &req); err != nil {
		h.logger.Errorf("ApproveSuggestion: invalid payload: %v", err)
		http.Error(w, `{"error":"Invalid suggestion payload"}`, http.StatusInternalServerError)
		return
	}

	if req.IconBase64 != "" {
		data, ext, err := decodeDataURL(req.IconBase64)
		if err != nil {
			h.logger.Errorf("ApproveSuggestion: decode icon: %v", err)
			http.Error(w, `{"error":"Invalid icon data"}`, http.StatusBadRequest)
			return
		}
		if err := os.MkdirAll("uploads/icons/weapons", 0755); err != nil {
			h.logger.Errorf("ApproveSuggestion: mkdir: %v", err)
			http.Error(w, `{"error":"Server error"}`, http.StatusInternalServerError)
			return
		}
		filename := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), id, ext)
		path := filepath.Join("uploads/icons/weapons", filename)
		if err := os.WriteFile(path, data, 0644); err != nil {
			h.logger.Errorf("ApproveSuggestion: write file: %v", err)
			http.Error(w, `{"error":"Server error"}`, http.StatusInternalServerError)
			return
		}
		req.Icon = fmt.Sprintf("/uploads/icons/weapons/%s", filename)
	}

	weaponID, err := h.weaponRepo.CreateWeapon(req)
	if err != nil {
		h.logger.Errorf("ApproveSuggestion: create weapon: %v", err)
		http.Error(w, `{"error":"Failed to create weapon"}`, http.StatusInternalServerError)
		return
	}

	if err := h.suggestionRepo.RemoveIconBase64(id); err != nil {
		h.logger.Warnf("ApproveSuggestion: remove base64: %v", err)
	}

	if err := h.suggestionRepo.UpdateStatus(id, claims.UserID, "approved"); err != nil {
		h.logger.Errorf("ApproveSuggestion: update status: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Suggestion approved",
		"weapon_id": weaponID,
	})
}

func (h *SuggestionHandler) Reject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error":"Invalid suggestion ID"}`, http.StatusBadRequest)
		return
	}

	claims := middleware.GetUserClaims(r)

	s, err := h.suggestionRepo.FindByID(id)
	if err != nil {
		h.logger.Errorf("RejectSuggestion: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	if s == nil {
		http.Error(w, `{"error":"Suggestion not found"}`, http.StatusNotFound)
		return
	}
	if s.Status != "pending" {
		http.Error(w, `{"error":"Suggestion already reviewed"}`, http.StatusBadRequest)
		return
	}

	if err := h.suggestionRepo.RemoveIconBase64(id); err != nil {
		h.logger.Warnf("RejectSuggestion: remove base64: %v", err)
	}

	if err := h.suggestionRepo.UpdateStatus(id, claims.UserID, "rejected"); err != nil {
		h.logger.Errorf("RejectSuggestion: update status: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Suggestion rejected"})
}

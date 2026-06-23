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
	"github.com/sirupsen/logrus"
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
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	var req models.CreateWeaponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(req.Description) > 500 {
		writeError(w, http.StatusBadRequest, "Description too long (max 500 chars)")
		return
	}

	raw, err := json.Marshal(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to encode payload")
		return
	}
	payload := json.RawMessage(raw)

	s, err := h.suggestionRepo.Create(claims.UserID, payload)
	if err != nil {
		h.logger.Errorf("CreateSuggestion: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
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
		writeError(w, http.StatusInternalServerError, "Database error")
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
		writeError(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}

	s, err := h.suggestionRepo.FindByID(id)
	if err != nil {
		h.logger.Errorf("GetSuggestion: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if s == nil {
		writeError(w, http.StatusNotFound, "Suggestion not found")
		return
	}

	json.NewEncoder(w).Encode(s)
}

func (h *SuggestionHandler) GetMy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	s, err := h.suggestionRepo.FindByID(id)
	if err != nil {
		h.logger.Errorf("GetMySuggestion: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if s == nil {
		writeError(w, http.StatusNotFound, "Suggestion not found")
		return
	}
	if s.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "Forbidden")
		return
	}

	json.NewEncoder(w).Encode(s)
}

func (h *SuggestionHandler) Approve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}

	claims := middleware.GetUserClaims(r)

	s, err := h.suggestionRepo.FindByID(id)
	if err != nil {
		h.logger.Errorf("ApproveSuggestion: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if s == nil {
		writeError(w, http.StatusNotFound, "Suggestion not found")
		return
	}
	if s.Status != "pending" {
		writeError(w, http.StatusBadRequest, "Suggestion already reviewed")
		return
	}

	var req models.CreateWeaponRequest
	if err := json.Unmarshal(s.Payload, &req); err != nil {
		h.logger.Errorf("ApproveSuggestion: invalid payload: %v", err)
		writeError(w, http.StatusInternalServerError, "Invalid suggestion payload")
		return
	}

	if req.IconBase64 != "" {
		data, ext, err := decodeDataURL(req.IconBase64)
		if err != nil {
			h.logger.Errorf("ApproveSuggestion: decode icon: %v", err)
			writeError(w, http.StatusBadRequest, "Invalid icon data")
			return
		}
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			h.logger.Errorf("ApproveSuggestion: mkdir: %v", err)
			writeError(w, http.StatusInternalServerError, "Server error")
			return
		}
		filename := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), id, ext)
		path := filepath.Join(uploadDir, filename)
		if err := os.WriteFile(path, data, 0644); err != nil {
			h.logger.Errorf("ApproveSuggestion: write file: %v", err)
			writeError(w, http.StatusInternalServerError, "Server error")
			return
		}
		req.Icon = fmt.Sprintf("/%s/%s", uploadDir, filename)
	}

	req.CreatedBy = &s.UserID

	weaponID, err := h.weaponRepo.CreateWeapon(req)
	if err != nil {
		h.logger.Errorf("ApproveSuggestion: create weapon: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to create weapon")
		return
	}

	if err := h.suggestionRepo.RemoveIconBase64(id); err != nil {
		h.logger.Warnf("ApproveSuggestion: remove base64: %v", err)
	}

	if err := h.suggestionRepo.UpdateStatus(id, claims.UserID, "approved", nil); err != nil {
		h.logger.Errorf("ApproveSuggestion: update status: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	h.logger.WithFields(logrus.Fields{
		"admin_id":  claims.UserID,
		"suggestion_id": id,
		"weapon_id": weaponID,
	}).Info("suggestion approved")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Suggestion approved",
		"weapon_id": weaponID,
	})
}

func (h *SuggestionHandler) Reject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}

	claims := middleware.GetUserClaims(r)

	s, err := h.suggestionRepo.FindByID(id)
	if err != nil {
		h.logger.Errorf("RejectSuggestion: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if s == nil {
		writeError(w, http.StatusNotFound, "Suggestion not found")
		return
	}
	if s.Status != "pending" {
		writeError(w, http.StatusBadRequest, "Suggestion already reviewed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	var body struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := h.suggestionRepo.RemoveIconBase64(id); err != nil {
		h.logger.Warnf("RejectSuggestion: remove base64: %v", err)
	}

	if err := h.suggestionRepo.UpdateStatus(id, claims.UserID, "rejected", &body.Reason); err != nil {
		h.logger.Errorf("RejectSuggestion: update status: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	h.logger.WithFields(logrus.Fields{
		"admin_id":  claims.UserID,
		"suggestion_id": id,
	}).Info("suggestion rejected")

	json.NewEncoder(w).Encode(map[string]string{"message": "Suggestion rejected"})
}

func (h *SuggestionHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	suggestions, err := h.suggestionRepo.FindByUserID(claims.UserID)
	if err != nil {
		h.logger.Errorf("ListMySuggestions: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if suggestions == nil {
		suggestions = []models.WeaponSuggestion{}
	}

	json.NewEncoder(w).Encode(suggestions)
}

func (h *SuggestionHandler) Resubmit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	s, err := h.suggestionRepo.FindByID(id)
	if err != nil {
		h.logger.Errorf("ResubmitSuggestion: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if s == nil {
		writeError(w, http.StatusNotFound, "Suggestion not found")
		return
	}
	if s.UserID != claims.UserID {
		writeError(w, http.StatusForbidden, "Not your suggestion")
		return
	}
	if s.Status != "rejected" {
		writeError(w, http.StatusBadRequest, "Only rejected suggestions can be resubmitted")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	var req models.CreateWeaponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(req.Description) > 500 {
		writeError(w, http.StatusBadRequest, "Description too long (max 500 chars)")
		return
	}

	raw, err := json.Marshal(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to encode payload")
		return
	}
	payload := json.RawMessage(raw)

	if err := h.suggestionRepo.UpdatePayload(id, payload); err != nil {
		h.logger.Errorf("ResubmitSuggestion: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":  claims.UserID,
		"suggestion_id": id,
	}).Info("suggestion resubmitted")

	updated, err := h.suggestionRepo.FindByID(id)
	if err != nil {
		h.logger.Errorf("ResubmitSuggestion: fetch updated: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (h *SuggestionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}

	if err := h.suggestionRepo.Delete(id); err != nil {
		h.logger.Errorf("DeleteSuggestion: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	h.logger.WithField("suggestion_id", id).Info("suggestion deleted")

	json.NewEncoder(w).Encode(map[string]string{"message": "Suggestion deleted"})
}

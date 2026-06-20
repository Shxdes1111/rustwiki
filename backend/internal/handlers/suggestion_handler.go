package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/logger"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/repository"
)

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

	var payload json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

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

	weaponID, err := h.weaponRepo.CreateWeapon(req)
	if err != nil {
		h.logger.Errorf("ApproveSuggestion: create weapon: %v", err)
		http.Error(w, `{"error":"Failed to create weapon"}`, http.StatusInternalServerError)
		return
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

	if err := h.suggestionRepo.UpdateStatus(id, claims.UserID, "rejected"); err != nil {
		h.logger.Errorf("RejectSuggestion: update status: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Suggestion rejected"})
}

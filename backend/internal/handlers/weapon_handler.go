package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/repository"
)

type WeaponHandler struct {
	weaponRepo repository.WeaponRepository
	Logger *logger.Logger
}

func NewWeaponHandler(weaponRepo repository.WeaponRepository, log *logger.Logger) *WeaponHandler {
	return &WeaponHandler{weaponRepo: weaponRepo, Logger: log}
}

func (h *WeaponHandler) GetWeapons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	weapons, err := h.weaponRepo.GetAllWeapons()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	json.NewEncoder(w).Encode(weapons)
}

func (h *WeaponHandler) GetWeapon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid weapon ID")
		return
	}

	weapon, err := h.weaponRepo.GetWeaponByID(id)
	if err != nil {
		h.Logger.Errorf("ERROR: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if weapon == nil {
		writeError(w, http.StatusNotFound, "Weapon not found")
		return
	}

	json.NewEncoder(w).Encode(weapon)
}

func (h *WeaponHandler) CreateWeapon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateWeaponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Errorf("CreateWeapon: invalid JSON: %v", err)
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	id, err := h.weaponRepo.CreateWeapon(req)
	if err != nil {
		h.Logger.Errorf("CreateWeapon: %v", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	h.Logger.WithField("weapon_id", id).Info("weapon created")

	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *WeaponHandler) DeleteWeapon(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid weapon ID")
		return
	}

	if err := h.weaponRepo.DeleteWeapon(id); err != nil {
		h.Logger.Errorf("DeleteWeapon: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	h.Logger.WithField("weapon_id", id).Info("weapon deleted")

	w.WriteHeader(http.StatusNoContent)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/logger"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	weapons, err := h.weaponRepo.GetAllWeapons()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(weapons)
}

func (h *WeaponHandler) GetWeapon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid weapon ID", http.StatusBadRequest)
		return
	}

	weapon, err := h.weaponRepo.GetWeaponByID(id)
	if err != nil {
		h.Logger.Errorf("ERROR: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if weapon == nil {
		http.Error(w, "Weapon not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(weapon)
}

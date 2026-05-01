package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/repository"
)

type WeaponHandler struct {
	weaponRepo repository.WeaponRepository
}

func NewWeaponHandler(weaponRepo repository.WeaponRepository) *WeaponHandler {
	return &WeaponHandler{weaponRepo: weaponRepo}
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

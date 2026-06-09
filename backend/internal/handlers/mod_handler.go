package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"backend/internal/logger"
	"backend/internal/repository"
)

type ModHandler struct {
	repo   repository.ModRepository
	Logger *logger.Logger
}

func NewModHandler(repo repository.ModRepository, log *logger.Logger) *ModHandler {
	return &ModHandler{repo: repo, Logger: log}
}

func (h *ModHandler) GetMod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid mod ID", http.StatusBadRequest)
		return
	}

	mod, err := h.repo.GetModByID(id)
	if err != nil {
		h.Logger.Errorf("ERROR: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if mod == nil {
		http.Error(w, "Mod not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(mod)
}

func (h *ModHandler) GetModList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	modList, err := h.repo.GetAllMods()
	if err != nil {
		h.Logger.Errorf("Database error while fetching mod list: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(modList)
}
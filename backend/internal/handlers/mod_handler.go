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
		writeError(w, http.StatusBadRequest, "Invalid mod ID")
		return
	}

	mod, err := h.repo.GetModByID(id)
	if err != nil {
		h.Logger.Errorf("ERROR: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if mod == nil {
		writeError(w, http.StatusNotFound, "Mod not found")
		return
	}

	json.NewEncoder(w).Encode(mod)
}

func (h *ModHandler) GetModList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	modList, err := h.repo.GetAllMods()
	if err != nil {
		h.Logger.Errorf("Database error while fetching mod list: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	json.NewEncoder(w).Encode(modList)
}
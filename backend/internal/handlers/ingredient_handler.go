package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/logger"
	"backend/internal/repository"
)

type IngredientHandler struct {
	repo   repository.IngredientRepository
	Logger *logger.Logger
}

func NewIngredientHandler(repo repository.IngredientRepository, log *logger.Logger) *IngredientHandler {
	return &IngredientHandler{repo: repo, Logger: log}
}

func (h *IngredientHandler) GetIngredientList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ingredients, err := h.repo.GetAllIngredients()
	if err != nil {
		h.Logger.Errorf("Database error while fetching ingredient list: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	json.NewEncoder(w).Encode(ingredients)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/repository"
)

// AmmoHandler отвечает за обработку HTTP-запросов, связанных с патронами
type AmmoHandler struct {
	repo   repository.AmmoRepository
	Logger *logger.Logger
}

// NewAmmoHandler создает новый экземпляр хэндлера патронов
func NewAmmoHandler(repo repository.AmmoRepository, log *logger.Logger) *AmmoHandler {
	return &AmmoHandler{repo: repo, Logger: log}
}

// GetAmmo обрабатывает запрос GET /api/ammo/{id}
func (h *AmmoHandler) GetAmmo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Извлекаем id из пути URL
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ammo ID")
		return
	}

	ammo, err := h.repo.GetAmmoByID(id)
	if err != nil {
		h.Logger.Errorf("Database error while fetching ammo: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Если патрон с таким id не найден
	if ammo == nil {
		writeError(w, http.StatusNotFound, "Ammo not found")
		return
	}

	// Отправляем JSON-ответ фронтенду
	json.NewEncoder(w).Encode(ammo)
}

func (h *AmmoHandler) GetAmmoList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ammoList, err := h.repo.GetAllAmmo()
	if err != nil {
		h.Logger.Errorf("Database error while fetching ammo list: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if ammoList == nil {
		ammoList = []models.Ammo{}
	}

	json.NewEncoder(w).Encode(ammoList)
}

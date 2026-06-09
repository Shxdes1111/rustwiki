package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"backend/internal/logger"
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
		http.Error(w, "Invalid ammo ID", http.StatusBadRequest)
		return
	}

	// Запрашиваем данные у репозитория
	ammo, err := h.repo.GetAmmoByID(id)
	if err != nil {
		h.Logger.Errorf("Database error while fetching ammo: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Если патрон с таким id не найден
	if ammo == nil {
		http.Error(w, "Ammo not found", http.StatusNotFound)
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
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ammoList)
}

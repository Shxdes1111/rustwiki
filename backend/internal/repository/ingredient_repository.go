package repository

import (
	"database/sql"

	"backend/internal/logger"
	"backend/internal/models"
)

type IngredientRepository interface {
	GetAllIngredients() ([]models.Ingredients, error)
}

type ingredientRepository struct {
	db  *sql.DB
	log *logger.Logger
}

func NewIngredientRepository(db *sql.DB, log *logger.Logger) IngredientRepository {
	return &ingredientRepository{db: db, log: log}
}

func (r *ingredientRepository) GetAllIngredients() ([]models.Ingredients, error) {
	r.log.Debug("GetAllIngredients: делаю запрос в таблицу ingredients")
	rows, err := r.db.Query("SELECT id, name, icon FROM ingredients ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Ingredients
	for rows.Next() {
		var ing models.Ingredients
		if err := rows.Scan(&ing.ID, &ing.Name, &ing.Icon); err != nil {
			return nil, err
		}
		list = append(list, ing)
	}

	return list, nil
}

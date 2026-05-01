package repository

import (
	"database/sql"

	"backend/internal/models"
)

type WeaponRepository interface {
	GetAllWeapons() ([]models.WeaponItem, error)
}

type weaponRepository struct {
	db *sql.DB
}

func NewWeaponRepository(db *sql.DB) WeaponRepository {
	return &weaponRepository{db: db}
}

func (r *weaponRepository) GetAllWeapons() ([]models.WeaponItem, error) {
	rows, err := r.db.Query("SELECT id, name, type FROM weapon_item")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weapons []models.WeaponItem
	for rows.Next() {
		var w models.WeaponItem
		if err := rows.Scan(&w.ID, &w.Name, &w.Type); err != nil {
			return nil, err
		}
		weapons = append(weapons, w)
	}

	return weapons, nil
}

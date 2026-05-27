package repository

import (
	"database/sql"

	"backend/internal/models"
)

type WeaponRepository interface {
	GetAllWeapons() ([]models.WeaponItem, error)
	GetWeaponByID(id int) (*models.WeaponItem, error)
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
		var weapon models.WeaponItem
		if err := rows.Scan(&weapon.ID, &weapon.Name, &weapon.Type); err != nil {
			return nil, err
		}
		weapons = append(weapons, weapon)
	}

	return weapons, nil
}

func (r *weaponRepository) GetWeaponByID(id int) (*models.WeaponItem, error) {
	row := r.db.QueryRow(
		"SELECT id, name, type, firemode, craftable, stacksize, category_id FROM weapon_item WHERE id = $1",
		id,
	)

	var weapon models.WeaponItem
	err := row.Scan(&weapon.ID, &weapon.Name, &weapon.Type, &weapon.Firemode, &weapon.Craftable, &weapon.Stacksize, &weapon.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// ammo
	ammoRows, err := r.db.Query("SELECT id, name, weapon_item_id FROM ammo WHERE weapon_item_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer ammoRows.Close()
	for ammoRows.Next() {
		var ammo models.Ammo
		if err := ammoRows.Scan(&ammo.ID, &ammo.Name, &ammo.WeaponItemID); err != nil {
			return nil, err
		}
		weapon.Ammo = append(weapon.Ammo, ammo)
	}

	// mods
	modRows, err := r.db.Query("SELECT id, name, weapon_item_id FROM mods WHERE weapon_item_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer modRows.Close()
	for modRows.Next() {
		var mods models.Mods
		if err := modRows.Scan(&mods.ID, &mods.Name, &mods.WeaponItemID); err != nil {
			return nil, err
		}
		weapon.Mods = append(weapon.Mods, mods)
	}

	// ingredients
	ingRows, err := r.db.Query("SELECT id, name, weapon_item_id FROM ingredients WHERE weapon_item_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer ingRows.Close()
	for ingRows.Next() {
		var ing models.Ingredients
		if err := ingRows.Scan(&ing.ID, &ing.Name, &ing.WeaponItemID); err != nil {
			return nil, err
		}
		weapon.Ingredients = append(weapon.Ingredients, ing)
	}

	return &weapon, nil
}

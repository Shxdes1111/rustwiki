package repository

import (
	"database/sql"

	"backend/internal/models"
)

type AmmoRepository interface {
	GetAmmoByID(id int) (*models.Ammo, error)
}

type ammoRepository struct {
	db *sql.DB
}

func NewAmmoRepository(db *sql.DB) AmmoRepository {
	return &ammoRepository{db: db}
}

func (r *ammoRepository) GetAmmoByID(id int) (*models.Ammo, error) {
	// 1. Получаем базовую информацию о патроне
	row := r.db.QueryRow(
		"SELECT id, name, icon, weapon_item_id FROM ammo WHERE id = $1",
		id,
	)

	var ammo models.Ammo
	err := row.Scan(&ammo.ID, &ammo.Name, &ammo.Icon, &ammo.WeaponItemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 2. Находим всё оружие, которое использует данный тип патронов
	// Делаем выборку по имени патрона, так как в Rust один калибр (например, 5.56) подходит к разным пушкам
	weaponRows, err := r.db.Query(`
		SELECT id, name, type, description, shortname, COALESCE(capacity, 0), COALESCE(time_to_craft, 0)
		FROM weapon_item 
		WHERE id IN (
			SELECT weapon_item_id FROM ammo WHERE id = $1 OR name = (SELECT name FROM ammo WHERE id = $1)
		)`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer weaponRows.Close()

	// 3. Наполняем слайс CompatibleWeapons внутри модели патрона
	for weaponRows.Next() {
		var weapon models.WeaponItem
		err := weaponRows.Scan(
			&weapon.ID, 
			&weapon.Name, 
			&weapon.Type, 
			&weapon.Description, 
			&weapon.Shortname, 
			&weapon.Capacity, 
			&weapon.TimeToCraft,
		)
		if err != nil {
			return nil, err
		}
		ammo.CompatibleWeapons = append(ammo.CompatibleWeapons, weapon)
	}

	return &ammo, nil
}

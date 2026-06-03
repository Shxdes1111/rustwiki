package repository

import (
	"database/sql"

	"backend/internal/logger"
	"backend/internal/models"
)

type WeaponRepository interface {
	GetAllWeapons() ([]models.WeaponItem, error)
	GetWeaponByID(id int) (*models.WeaponItem, error)
}

type weaponRepository struct {
	db  *sql.DB
	log *logger.Logger
}

func NewWeaponRepository(db *sql.DB, log *logger.Logger) WeaponRepository {
	return &weaponRepository{db: db, log: log}
}

func (r *weaponRepository) GetAllWeapons() ([]models.WeaponItem, error) {
	r.log.Info("GetAllWeapons: делаю запрос в таблицу weapon_item")
	rows, err := r.db.Query("SELECT id, name, type, description, shortname, COALESCE(capacity, 0), COALESCE(time_to_craft, 0) FROM weapon_item")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weapons []models.WeaponItem
	for rows.Next() {
		var weapon models.WeaponItem
		if err := rows.Scan(&weapon.ID, &weapon.Name, &weapon.Type, &weapon.Description, &weapon.Shortname, &weapon.Capacity, &weapon.TimeToCraft); err != nil {
			return nil, err
		}
		weapons = append(weapons, weapon)
	}

	return weapons, nil
}

func (r *weaponRepository) GetWeaponByID(id int) (*models.WeaponItem, error) {
	r.log.Info("GetWeaponByID: делаю запрос в таблицу weapon_item")
	row := r.db.QueryRow(
		"SELECT id, name, type, firemode, craftable, stacksize, description, shortname, COALESCE(capacity, 0), COALESCE(time_to_craft, 0), category_id FROM weapon_item WHERE id = $1",
		id,
	)

	var weapon models.WeaponItem
	err := row.Scan(&weapon.ID, &weapon.Name, &weapon.Type, &weapon.Firemode, &weapon.Craftable, &weapon.Stacksize, &weapon.Description, &weapon.Shortname, &weapon.Capacity, &weapon.TimeToCraft, &weapon.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// ammo
	r.log.Info("GetWeaponByID: делаю запрос в таблицы ammo и weapon_ammo")
	ammoRows, err := r.db.Query(`
        SELECT a.id, a.name, a.icon 
        FROM ammo a
        JOIN weapon_ammo wa ON a.id = wa.ammo_id
        WHERE wa.weapon_item_id = $1`, 
        id,
    )
    if err != nil {
        return nil, err
    }
    defer ammoRows.Close()

    for ammoRows.Next() {
        var ammo models.Ammo
        if err := ammoRows.Scan(&ammo.ID, &ammo.Name, &ammo.Icon); err != nil {
            return nil, err
        }
        weapon.Ammo = append(weapon.Ammo, ammo)
    }

	// mods
	r.log.Info("GetWeaponByID: делаю запрос в таблицы mods и weapon_mods")
	modRows, err := r.db.Query(`
		SELECT m.id, m.name, m.icon, wm.weapon_item_id 
		FROM mods m
		JOIN weapon_mods wm ON m.id = wm.mod_id
		WHERE wm.weapon_item_id = $1`, 
		id,
	)
	if err != nil {
		return nil, err
	}
	defer modRows.Close()

	for modRows.Next() {
		var mod models.Mods
		// Теперь данные считываются корректно, включая id оружия из связующей таблицы
		if err := modRows.Scan(&mod.ID, &mod.Name, &mod.Icon, &mod.WeaponItemID); err != nil {
			return nil, err
		}
		weapon.Mods = append(weapon.Mods, mod)
	}

	// ingredients
	r.log.Info("GetWeaponByID: делаю запрос в таблицы ingredients и weapon_ingredients")
	ingRows, err := r.db.Query(`
		SELECT i.id, i.name, wi.amount, i.icon 
		FROM ingredients i
		JOIN weapon_ingredients wi ON i.id = wi.ingredients_id
		WHERE wi.weapon_item_id = $1`, 
		id,
	)
	if err != nil {
		return nil, err
	}
	defer ingRows.Close()

	for ingRows.Next() {
		var ing models.Ingredients
		// Переменные сканируются в том же порядке, сохраняя структуру models.Ingredients без изменений
		if err := ingRows.Scan(&ing.ID, &ing.Name, &ing.Amount, &ing.Icon); err != nil {
			return nil, err
		}
		ing.WeaponItemID = &id 
		weapon.Ingredients = append(weapon.Ingredients, ing)
	}

	return &weapon, nil
}

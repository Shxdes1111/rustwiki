package repository

import (
	"database/sql"

	"backend/internal/logger"
	"backend/internal/models"
)

type ModRepository interface {
	GetModByID(id int) (*models.Mods, error)
	GetAllMods() ([]models.Mods, error)
}

type modRepository struct {
	db  *sql.DB
	log *logger.Logger
}

func NewModRepository(db *sql.DB, log *logger.Logger) ModRepository {
	return &modRepository{db: db, log: log}
}

func (r *modRepository) GetModByID(id int) (*models.Mods, error) {
	// 1. Получаем сам модуль
	r.log.Debug("GetModByID: делаю запрос в таблицу mods")
	row := r.db.QueryRow(
		"SELECT id, name, icon FROM mods WHERE id = $1",
		id,
	)

	var mod models.Mods
	err := row.Scan(&mod.ID, &mod.Name, &mod.Icon)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 2. Находим всё оружие, к которому подходит этот модуль (Many-to-Many)
	// Используем честный JOIN через таблицу weapon_mods
	r.log.Debug("GetModByID: делаю запрос в таблицы weapon_item и weapon_mods")
	weaponRows, err := r.db.Query(`
		SELECT w.id, w.name, w.type, w.description, w.shortname, w.icon, COALESCE(w.capacity, 0), COALESCE(w.time_to_craft, 0)
		FROM weapon_item w
		JOIN weapon_mods wm ON w.id = wm.weapon_item_id
		WHERE wm.mod_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer weaponRows.Close()

	// 3. Сканируем оружие в слайс CompatibleWeapons внутри модели модуля
	// (Убедитесь, что в структуре models.Mods у вас есть поле CompatibleWeapons []WeaponItem или аналогичное)
	for weaponRows.Next() {
		var weapon models.WeaponItem
		err := weaponRows.Scan(
			&weapon.ID, 
			&weapon.Name, 
			&weapon.Type, 
			&weapon.Description, 
			&weapon.Shortname, 
			&weapon.Icon, 
			&weapon.Capacity, 
			&weapon.TimeToCraft,
		)
		if err != nil {
			return nil, err
		}
		mod.CompatibleWeapons = append(mod.CompatibleWeapons, weapon)
	}

	return &mod, nil
}

func (r *modRepository) GetAllMods() ([]models.Mods, error) {
	r.log.Debug("GetAllMods: делаю запрос в таблицу mods")
	rows, err := r.db.Query("SELECT id, name, icon FROM mods ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modList []models.Mods
	for rows.Next() {
		var m models.Mods
		if err := rows.Scan(&m.ID, &m.Name, &m.Icon); err != nil {
			return nil, err
		}
		modList = append(modList, m)
	}

	return modList, nil
}

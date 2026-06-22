package repository

import (
	"database/sql"

	"backend/internal/logger"
	"backend/internal/models"
)

type AmmoRepository interface {
	GetAmmoByID(id int) (*models.Ammo, error)
	GetAllAmmo() ([]models.Ammo, error)
}

type ammoRepository struct {
	db  *sql.DB
	log *logger.Logger
}

func NewAmmoRepository(db *sql.DB, log *logger.Logger) AmmoRepository {
	return &ammoRepository{db: db, log: log}
}

func (r *ammoRepository) GetAmmoByID (id int) (*models.Ammo, error) {
    // 1. Получаем базовую информацию о патроне (убрали weapon_item_id из SELECT)
    r.log.Info("GetAmmoByID: делаю запрос в таблицу ammo")
    row := r.db.QueryRow(
        "SELECT id, name, icon FROM ammo WHERE id = $1",
        id,
    )
    
    var ammo models.Ammo
    err := row.Scan(&ammo.ID, &ammo.Name, &ammo.Icon)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    // 2. Находим всё оружие через связующую таблицу weapon_ammo (Честный JOIN)
    r.log.Info("GetAmmoByID: делаю запрос в таблицы weapon_item и weapon_ammo")
    weaponRows, err := r.db.Query(`
        SELECT w.id, w.name, w.type, w.description, w.shortname, w.icon,
               COALESCE(w.capacity, 0), COALESCE(w.time_to_craft, 0)
        FROM weapon_item w
        JOIN weapon_ammo wa ON w.id = wa.weapon_item_id
        WHERE wa.ammo_id = $1`, 
        id,
    )
    if err != nil {
        return nil, err
    }
    defer weaponRows.Close()

    // 3. Наполняем слайс внутри патрона
    for weaponRows.Next() {
        var weapon models.WeaponItem
        err := weaponRows.Scan(
            &weapon.ID, &weapon.Name, &weapon.Type, &weapon.Description,
            &weapon.Shortname, &weapon.Icon, &weapon.Capacity, &weapon.TimeToCraft,
        )
        if err != nil {
            return nil, err
        }
        ammo.CompatibleWeapons = append(ammo.CompatibleWeapons, weapon)
	}
	
	return &ammo, nil
}

func (r *ammoRepository) GetAllAmmo() ([]models.Ammo, error) {
	r.log.Debug("GetAllAmmo: делаю запрос в таблицу ammo")
	rows, err := r.db.Query("SELECT id, name, icon FROM ammo ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ammoList []models.Ammo
	for rows.Next() {
		var a models.Ammo
		if err := rows.Scan(&a.ID, &a.Name, &a.Icon); err != nil {
			return nil, err
		}
		ammoList = append(ammoList, a)
	}

	return ammoList, nil
}

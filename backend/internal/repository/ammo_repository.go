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

func (r *ammoRepository) GetAmmoByID (id int) (*models.Ammo, error) {
    // 1. Получаем базовую информацию о патроне (убрали weapon_item_id из SELECT)
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
    weaponRows, err := r.db.Query(`
        SELECT w.id, w.name, w.type, w.description, w.shortname, 
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
            &weapon.Shortname, &weapon.Capacity, &weapon.TimeToCraft,
        )
        if err != nil {
            return nil, err
        }
        ammo.CompatibleWeapons = append(ammo.CompatibleWeapons, weapon)
    }
    
    return &ammo, nil
}

package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"backend/internal/logger"
	"backend/internal/models"
)

type WeaponRepository interface {
	GetAllWeapons() ([]models.WeaponItem, error)
	GetWeaponByID(id int) (*models.WeaponItem, error)
	CreateWeapon(req models.CreateWeaponRequest) (int, error)
	DeleteWeapon(id int) error
	FindByUserID(userID int) ([]models.WeaponItem, error)
}

type weaponRepository struct {
	db        *sql.DB
	log       *logger.Logger
	publicURL string
}

func NewWeaponRepository(db *sql.DB, log *logger.Logger, publicURL string) WeaponRepository {
	return &weaponRepository{db: db, log: log, publicURL: publicURL}
}

func (r *weaponRepository) GetAllWeapons() ([]models.WeaponItem, error) {
	r.log.Debug("GetAllWeapons: делаю запрос в таблицу weapon_item")
	rows, err := r.db.Query("SELECT id, name, type, firemode, craftable, stacksize, description, shortname, icon, capacity, time_to_craft, category_id FROM weapon_item ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weapons []models.WeaponItem
	for rows.Next() {
		var weapon models.WeaponItem
		var cap, ttc sql.NullInt64
		if err := rows.Scan(&weapon.ID, &weapon.Name, &weapon.Type, &weapon.Firemode, &weapon.Craftable, &weapon.Stacksize, &weapon.Description, &weapon.Shortname, &weapon.Icon, &cap, &ttc, &weapon.CategoryID); err != nil {
			return nil, err
		}
		if cap.Valid { v := int(cap.Int64); weapon.Capacity = &v }
		if ttc.Valid { v := int(ttc.Int64); weapon.TimeToCraft = &v }
		weapons = append(weapons, weapon)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return weapons, nil
}

func (r *weaponRepository) GetWeaponByID(id int) (*models.WeaponItem, error) {
	r.log.Debug("GetWeaponByID: делаю запрос в таблицу weapon_item")

	var weapon models.WeaponItem
	var cap, ttc sql.NullInt64
	var createdBy sql.NullInt64
	err := r.db.QueryRow(`
		SELECT w.id, w.name, w.type, w.firemode, w.craftable, w.stacksize,
		       w.description, w.shortname, w.icon, w.capacity, w.time_to_craft,
		       w.category_id, w.views, w.created_by, COALESCE(u.username, '')
		FROM weapon_item w
		LEFT JOIN users u ON w.created_by = u.id
		WHERE w.id = $1`,
		id,
	).Scan(&weapon.ID, &weapon.Name, &weapon.Type, &weapon.Firemode, &weapon.Craftable,
		&weapon.Stacksize, &weapon.Description, &weapon.Shortname, &weapon.Icon,
		&cap, &ttc, &weapon.CategoryID, &weapon.Views, &createdBy, &weapon.AuthorName)
	if cap.Valid {
		v := int(cap.Int64)
		weapon.Capacity = &v
	}
	if ttc.Valid {
		v := int(ttc.Int64)
		weapon.TimeToCraft = &v
	}
	if createdBy.Valid {
		v := int(createdBy.Int64)
		weapon.CreatedBy = &v
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// increment views atomically
	if err := r.db.QueryRow(`UPDATE weapon_item SET views = views + 1 WHERE id = $1 RETURNING views`, id).Scan(&weapon.Views); err != nil {
		r.log.Warnf("GetWeaponByID: failed to increment views: %v", err)
	}

	// ammo
	r.log.Debug("GetWeaponByID: делаю запрос в таблицы ammo и weapon_ammo")
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
    if err := ammoRows.Err(); err != nil {
        return nil, err
    }

	// mods
	r.log.Debug("GetWeaponByID: делаю запрос в таблицы mods и weapon_mods")
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
	if err := modRows.Err(); err != nil {
		return nil, err
	}

	// ingredients
	r.log.Debug("GetWeaponByID: делаю запрос в таблицы ingredients и weapon_ingredients")
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
	if err := ingRows.Err(); err != nil {
		return nil, err
	}

	return &weapon, nil
}

func (r *weaponRepository) CreateWeapon(req models.CreateWeaponRequest) (int, error) {
	if len(req.Description) > 500 {
		return 0, fmt.Errorf("description too long (max 500 chars)")
	}
	r.log.Debug("CreateWeapon: начинаю транзакцию создания оружия")

	var err error
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var newID int
	err = tx.QueryRow(`
		INSERT INTO weapon_item (name, type, firemode, craftable, stacksize, description, shortname, icon, capacity, time_to_craft, category_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`,
		req.Name, req.Type, req.Firemode, req.Craftable, req.Stacksize,
		req.Description, req.Shortname, req.Icon, req.Capacity, req.TimeToCraft, req.CategoryID,
		req.CreatedBy,
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	for _, ammoID := range req.AmmoIDs {
		if _, err := tx.Exec(`INSERT INTO weapon_ammo (weapon_item_id, ammo_id) VALUES ($1, $2)`, newID, ammoID); err != nil {
			return 0, err
		}
	}

	for _, modID := range req.ModIDs {
		if _, err := tx.Exec(`INSERT INTO weapon_mods (weapon_item_id, mod_id) VALUES ($1, $2)`, newID, modID); err != nil {
			return 0, err
		}
	}

	for _, ing := range req.Ingredients {
		if _, err := tx.Exec(`INSERT INTO weapon_ingredients (weapon_item_id, ingredients_id, amount) VALUES ($1, $2, $3)`, newID, ing.ID, ing.Amount); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	r.log.Debugf("CreateWeapon: оружие создано с id=%d", newID)
	return newID, nil
}

func (r *weaponRepository) FindByUserID(userID int) ([]models.WeaponItem, error) {
	rows, err := r.db.Query(`
		SELECT w.id, w.name, w.type, w.firemode, w.craftable, w.stacksize,
		       w.description, w.shortname, w.icon, w.capacity, w.time_to_craft,
		       w.category_id, w.views, w.created_by, COALESCE(u.username, '')
		FROM weapon_item w
		LEFT JOIN users u ON w.created_by = u.id
		WHERE w.created_by = $1
		ORDER BY w.id ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weapons []models.WeaponItem
	for rows.Next() {
		var weapon models.WeaponItem
		var cap, ttc sql.NullInt64
		var createdBy sql.NullInt64
		if err := rows.Scan(&weapon.ID, &weapon.Name, &weapon.Type, &weapon.Firemode,
			&weapon.Craftable, &weapon.Stacksize, &weapon.Description, &weapon.Shortname,
			&weapon.Icon, &cap, &ttc, &weapon.CategoryID, &weapon.Views, &createdBy,
			&weapon.AuthorName); err != nil {
			return nil, err
		}
		if cap.Valid {
			v := int(cap.Int64)
			weapon.Capacity = &v
		}
		if ttc.Valid {
			v := int(ttc.Int64)
			weapon.TimeToCraft = &v
		}
		if createdBy.Valid {
			v := int(createdBy.Int64)
			weapon.CreatedBy = &v
		}
		weapons = append(weapons, weapon)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return weapons, nil
}

func (r *weaponRepository) DeleteWeapon(id int) error {
	r.log.Debugf("DeleteWeapon: удаляю оружие id=%d", id)

	var icon *string
	err := r.db.QueryRow("SELECT icon FROM weapon_item WHERE id = $1", id).Scan(&icon)
	if err != nil {
		return err
	}

	if icon != nil {
		cleanPath := *icon
		if strings.HasPrefix(cleanPath, r.publicURL+"/uploads/") {
			cleanPath = strings.TrimPrefix(cleanPath, r.publicURL)
		}
		if strings.HasPrefix(cleanPath, "/uploads/") {
			relPath := strings.TrimPrefix(cleanPath, "/uploads/")
			absPath := filepath.Join("uploads", relPath)
			if err := os.Remove(absPath); err != nil && !os.IsNotExist(err) {
				r.log.Warnf("DeleteWeapon: не удалось удалить файл %s: %v", absPath, err)
			} else {
				r.log.Debugf("DeleteWeapon: удалён файл %s", absPath)
			}
		}
	}

	_, err = r.db.Exec("DELETE FROM weapon_item WHERE id = $1", id)
	return err
}

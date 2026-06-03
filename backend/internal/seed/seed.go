package seed

import (
	"database/sql"
	"fmt"

	"backend/internal/logger"
)

type weaponSeed struct {
	ID          int
	Name        string
	Type        string
	Firemode    string
	Craftable   bool
	Stacksize   int
	Description string
	Shortname   string
	Capacity    *int
	TimeToCraft *int
	CategoryID  int
	AmmoIDs     []int
	ModIDs      []int
	Ingredients []ingredientSeed
}

type ingredientSeed struct {
	ID     int
	Amount int
}

func ptr(i int) *int { return &i }

func Seed(db *sql.DB, log *logger.Logger) error {
	log.Info("Seeding database...")

	if _, err := db.Exec(`INSERT INTO category (id, name) VALUES (1, 'weapons'), (2, 'armor') ON CONFLICT (id) DO NOTHING`); err != nil {
		return fmt.Errorf("category: %w", err)
	}

	ammoList := []struct {
		ID   int
		Name string
		Icon string
	}{
		{1, "5.56mm Rifle Ammo", "/icons/ammo/rifle.png"},
		{2, "HV 5.56mm Rifle Ammo", "/icons/ammo/rifle.png"},
		{3, "Incendiary 5.56mm Rifle Ammo", "/icons/ammo/rifle.png"},
		{4, "Explosive 5.56mm Rifle Ammo", "/icons/ammo/rifle.png"},
		{5, "9mm Pistol Ammo", "/icons/ammo/pistol.png"},
		{6, "HV 9mm Pistol Ammo", "/icons/ammo/pistol.png"},
		{7, "Incendiary 9mm Pistol Ammo", "/icons/ammo/pistol.png"},
		{8, "Handmade Shell", "/icons/ammo/shells.png"},
		{9, "12 Gauge Buckshot", "/icons/ammo/shells.png"},
		{10, "12 Gauge Slug", "/icons/ammo/shells.png"},
		{11, "12 Gauge Incendiary Shell", "/icons/ammo/shells.png"},
		{12, "Wooden Arrow", "/icons/ammo/arrow.png"},
		{13, "High Velocity Arrow", "/icons/ammo/arrow.png"},
		{14, "Bone Arrow", "/icons/ammo/arrow.png"},
		{15, "Fire Arrow", "/icons/ammo/arrow.png"},
	}
	for _, a := range ammoList {
		if _, err := db.Exec(`INSERT INTO ammo (id, name, icon) VALUES ($1,$2,$3) ON CONFLICT (id) DO UPDATE SET name=$2, icon=$3`, a.ID, a.Name, a.Icon); err != nil {
			return fmt.Errorf("ammo %d: %w", a.ID, err)
		}
	}

	modList := []struct {
		ID   int
		Name string
		Icon string
	}{
		{1, "Weapon flashlight", "Weapon flashlight"},
		{2, "Weapon Lasersight", "/icons/mods/Weapon Lasersight.png"},
		{3, "Soda Can Silencer", "/icons/mods/Soda Can Silencer.png"},
		{4, "Oil Filter Silencer", "/icons/mods/Oil Filter Silencer.png"},
		{5, "Military Silencer", "/icons/mods/silencer.png"},
		{6, "Simple Handmade Sight", "Simple Handmade Sight"},
		{7, "Holosight", "/icons/mods/holo.png"},
		{8, "8x Zoom Scope", "/icons/mods/scope.png"},
		{9, "Variable Zoom Scope", "Variable Zoom Scope"},
		{10, "Gas Compression Overdrive", "Gas Compression Overdrive"},
		{11, "Muzzle Boost", "/icons/mods/muzzle.png"},
		{12, "Muzzle Brake", "/icons/mods/muzzle.png"},
		{13, "Burst Module", "Burst Module"},
		{14, "Extended Magazine", "Extended Magazine"},
		{15, "Targeting Attachment", "Targeting Attachment"},
	}
	for _, m := range modList {
		if _, err := db.Exec(`INSERT INTO mods (id, name, icon) VALUES ($1,$2,$3) ON CONFLICT (id) DO UPDATE SET name=$2, icon=$3`, m.ID, m.Name, m.Icon); err != nil {
			return fmt.Errorf("mods %d: %w", m.ID, err)
		}
	}

	ingredientList := []struct {
		ID   int
		Name string
		Icon string
	}{
		{1, "Metal Fragments", "/icons/ingredients/metal_frags.png"},
		{2, "Wood", "/icons/ingredients/wood.png"},
		{3, "Cloth", "/icons/ingredients/cloth.png"},
		{4, "High Quality Metal", "/icons/ingredients/hqm.png"},
		{5, "Animal Fat", "/icons/ingredients/animal_fat.png"},
	}
	for _, i := range ingredientList {
		if _, err := db.Exec(`INSERT INTO ingredients (id, name, icon) VALUES ($1,$2,$3) ON CONFLICT (id) DO UPDATE SET name=$2, icon=$3`, i.ID, i.Name, i.Icon); err != nil {
			return fmt.Errorf("ingredients %d: %w", i.ID, err)
		}
	}

	clothingList := []struct {
		ID            int
		Name          string
		Equipmentslot string
		Protection    int
		Craftable     bool
		Stacksize     int
		CategoryID    int
	}{
		{1, "Vest", "body", 70, true, 1, 2},
		{2, "Hoodie", "body", 35, true, 1, 2},
		{3, "Jacket", "body", 45, true, 1, 2},
		{4, "T-Shirt", "body", 0, true, 1, 2},
		{5, "Jeans", "legs", 35, true, 1, 2},
		{6, "Boots", "feet", 10, true, 1, 2},
		{7, "Baseball Cap", "head", 0, true, 1, 2},
		{8, "Balaclava", "head", 5, true, 1, 2},
		{9, "Bone Helmet", "head", 15, true, 1, 2},
		{10, "Coffee Can Helmet", "head", 20, true, 1, 2},
	}
	for _, c := range clothingList {
		if _, err := db.Exec(`INSERT INTO clothing_item (id, name, equipmentslot, protection, craftable, stacksize, category_id) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (id) DO UPDATE SET name=$2, equipmentslot=$3, protection=$4, craftable=$5, stacksize=$6, category_id=$7`,
			c.ID, c.Name, c.Equipmentslot, c.Protection, c.Craftable, c.Stacksize, c.CategoryID); err != nil {
			return fmt.Errorf("clothing_item %d: %w", c.ID, err)
		}
	}

	weapons := []weaponSeed{
		{
			ID: 1, Name: "Assault Rifle", Type: "range", Firemode: "automatic", Craftable: true, Stacksize: 1,
			Description: "The Assault Rifle is a powerful automatic rifle, deadly at medium range. A favorite among Rust players.",
			Shortname: "ar", Capacity: ptr(30), TimeToCraft: ptr(15), CategoryID: 1,
			AmmoIDs: []int{1, 2, 3, 4},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 14, 15},
			Ingredients: []ingredientSeed{{1, 250}},
		},
		{
			ID: 2, Name: "LR-300", Type: "range", Firemode: "automatic", Craftable: false, Stacksize: 1,
			Description: "The LR-300 is a high-accuracy automatic rifle. Cannot be crafted — only found in loot crates.",
			Shortname: "lr300", Capacity: ptr(30), TimeToCraft: nil, CategoryID: 1,
			AmmoIDs: []int{1, 2, 3, 4},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 14, 15},
			Ingredients: []ingredientSeed{{1, 150}},
		},
		{
			ID: 3, Name: "MP5A4", Type: "range", Firemode: "automatic", Craftable: true, Stacksize: 1,
			Description: "The MP5A4 is a compact submachine gun with a high rate of fire. Effective in close quarters.",
			Shortname: "mp5a4", Capacity: ptr(30), TimeToCraft: ptr(10), CategoryID: 1,
			AmmoIDs: []int{5, 6, 7},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 14},
			Ingredients: []ingredientSeed{{2, 200}},
		},
		{
			ID: 4, Name: "Thompson", Type: "range", Firemode: "automatic", Craftable: true, Stacksize: 1,
			Description: "The Thompson submachine gun, also known as the \"Tommy Gun\", is a classic automatic weapon.",
			Shortname: "tommy", Capacity: ptr(20), TimeToCraft: ptr(10), CategoryID: 1,
			AmmoIDs: []int{5, 6, 7},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		{
			ID: 5, Name: "Pump Shotgun", Type: "range", Firemode: "semi", Craftable: true, Stacksize: 1,
			Description: "A pump-action shotgun with a tight spread. Devastating at close range.",
			Shortname: "pump_shotgun", Capacity: ptr(6), TimeToCraft: ptr(15), CategoryID: 1,
			AmmoIDs: []int{8, 9, 10, 11},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12},
		},
		{
			ID: 6, Name: "Double Barrel Shotgun", Type: "range", Firemode: "double", Craftable: true, Stacksize: 1,
			Description: "A double-barreled shotgun. Two shots, high damage, slow reload.",
			Shortname: "double_barrel", Capacity: ptr(2), TimeToCraft: ptr(10), CategoryID: 1,
			AmmoIDs: []int{8, 9, 10, 11},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12},
		},
		{
			ID: 7, Name: "Revolver", Type: "range", Firemode: "semi", Craftable: true, Stacksize: 1,
			Description: "A six-shot revolver. Slower fire rate but packs a punch.",
			Shortname: "revolver", Capacity: ptr(6), TimeToCraft: ptr(5), CategoryID: 1,
			AmmoIDs: []int{5, 6, 7},
			ModIDs:  []int{3, 4, 5, 11, 12},
		},
		{
			ID: 8, Name: "Semi-Automatic Rifle", Type: "range", Firemode: "semi", Craftable: true, Stacksize: 1,
			Description: "A semi-automatic pistol. Reliable and easy to craft.",
			Shortname: "semi_auto_pistol", Capacity: ptr(8), TimeToCraft: ptr(5), CategoryID: 1,
			AmmoIDs: []int{1, 2, 3, 4},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 15},
		},
		{
			ID: 9, Name: "Machete", Type: "melee", Firemode: "none", Craftable: true, Stacksize: 1,
			Description: "A sharp machete. No ammo needed, just swing.",
			Shortname: "machete", Capacity: nil, TimeToCraft: ptr(5), CategoryID: 1,
			AmmoIDs: []int{},
			ModIDs:  []int{},
		},
		{
			ID: 10, Name: "Hunting Bow", Type: "range", Firemode: "none", Craftable: true, Stacksize: 1,
			Description: "A simple hunting bow. Silent but deadly with practice.",
			Shortname: "hunting_bow", Capacity: ptr(1), TimeToCraft: ptr(5), CategoryID: 1,
			AmmoIDs: []int{12, 13, 14, 15},
			ModIDs:  []int{},
		},
		{
			ID: 11, Name: "Compound Bow", Type: "range", Firemode: "none", Craftable: true, Stacksize: 1,
			Description: "A modern compound bow. Higher damage and faster arrow velocity than the standard bow.",
			Shortname: "compound_bow", Capacity: ptr(1), TimeToCraft: ptr(10), CategoryID: 1,
			AmmoIDs: []int{12, 13, 14, 15},
			ModIDs:  []int{},
		},
	}

	for _, w := range weapons {
		if _, err := db.Exec(`
			INSERT INTO weapon_item (id, name, type, firemode, craftable, stacksize, description, shortname, capacity, time_to_craft, category_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
			ON CONFLICT (id) DO UPDATE SET name=$2, type=$3, firemode=$4, craftable=$5, stacksize=$6, description=$7, shortname=$8, capacity=$9, time_to_craft=$10, category_id=$11`,
			w.ID, w.Name, w.Type, w.Firemode, w.Craftable, w.Stacksize, w.Description, w.Shortname, w.Capacity, w.TimeToCraft, w.CategoryID,
		); err != nil {
			return fmt.Errorf("weapon_item %d: %w", w.ID, err)
		}

		if _, err := db.Exec(`DELETE FROM weapon_ammo WHERE weapon_item_id = $1`, w.ID); err != nil {
			return fmt.Errorf("weapon_ammo delete %d: %w", w.ID, err)
		}
		for _, ammoID := range w.AmmoIDs {
			if _, err := db.Exec(`INSERT INTO weapon_ammo (weapon_item_id, ammo_id) VALUES ($1, $2)`, w.ID, ammoID); err != nil {
				return fmt.Errorf("weapon_ammo %d/%d: %w", w.ID, ammoID, err)
			}
		}

		if _, err := db.Exec(`DELETE FROM weapon_mods WHERE weapon_item_id = $1`, w.ID); err != nil {
			return fmt.Errorf("weapon_mods delete %d: %w", w.ID, err)
		}
		for _, modID := range w.ModIDs {
			if _, err := db.Exec(`INSERT INTO weapon_mods (weapon_item_id, mod_id) VALUES ($1, $2)`, w.ID, modID); err != nil {
				return fmt.Errorf("weapon_mods %d/%d: %w", w.ID, modID, err)
			}
		}

		if _, err := db.Exec(`DELETE FROM weapon_ingredients WHERE weapon_item_id = $1`, w.ID); err != nil {
			return fmt.Errorf("weapon_ingredients delete %d: %w", w.ID, err)
		}
		for _, ing := range w.Ingredients {
			if _, err := db.Exec(`INSERT INTO weapon_ingredients (weapon_item_id, ingredients_id, amount) VALUES ($1, $2, $3)`, w.ID, ing.ID, ing.Amount); err != nil {
				return fmt.Errorf("weapon_ingredients %d/%d: %w", w.ID, ing.ID, err)
			}
		}
	}

	log.Info("Database seeded successfully")
	return nil
}

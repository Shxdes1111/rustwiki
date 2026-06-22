package seed

import (
	"database/sql"
	"fmt"

	"backend/internal/logger"

	"golang.org/x/crypto/bcrypt"
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
	Icon        string
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

	admin1Hash, err := bcrypt.GenerateFromPassword([]byte("926754"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("admin password hash: %w", err)
	}
	if _, err := db.Exec(`INSERT INTO users (username, password_hash, role) VALUES ('shxdes', $1, 'admin') ON CONFLICT (username) DO UPDATE SET password_hash = $1, role = 'admin'`, string(admin1Hash)); err != nil {
		return fmt.Errorf("shxdes user: %w", err)
	}

	if _, err := db.Exec(`INSERT INTO category (id, name) VALUES (1, 'weapons'), (2, 'armor') ON CONFLICT (id) DO NOTHING`); err != nil {
		return fmt.Errorf("category: %w", err)
	}

	ammoList := []struct {
		ID   int
		Name string
		Icon string
	}{
		{1, "5.56mm Rifle Ammo", "/icons/ammo/ammo-rifle.avif"},
		{2, "HV 5.56mm Rifle Ammo", "/icons/ammo/ammo-rifle-hv.avif"},
		{3, "Incendiary 5.56mm Rifle Ammo", "/icons/ammo/ammo-rifle-incendiary.avif"},
		{4, "Explosive 5.56mm Rifle Ammo", "/icons/ammo/ammo-rifle-explosive.avif"},
		{5, "9mm Pistol Ammo", "/icons/ammo/ammo-pistol.avif"},
		{6, "HV 9mm Pistol Ammo", "/icons/ammo/ammo-pistol-hv.avif"},
		{7, "Incendiary 9mm Pistol Ammo", "/icons/ammo/ammo-pistol-fire.avif"},
		{8, "Handmade Shell", "/icons/ammo/ammo-handmade-shell.avif"},
		{9, "12 Gauge Buckshot", "/icons/ammo/ammo-shotgun.avif"},
		{10, "12 Gauge Slug", "/icons/ammo/ammo-shotgun-slug.avif"},
		{11, "12 Gauge Incendiary Shell", "/icons/ammo/ammo-shotgun-fire.avif"},
		{12, "Wooden Arrow", "/icons/ammo/arrow-wooden.avif"},
		{13, "High Velocity Arrow", "/icons/ammo/arrow-hv.avif"},
		{14, "Bone Arrow", "/icons/ammo/arrow-bone.avif"},
		{15, "Fire Arrow", "/icons/ammo/arrow-fire.avif"},
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
		{1, "Weapon flashlight", "/icons/mods/weapon-mod-flashlight.avif"},
		{2, "Weapon Lasersight", "/icons/mods/weapon-mod-lasersight.avif"},
		{3, "Soda Can Silencer", "/icons/mods/weapon-mod-sodacansilencer.avif"},
		{4, "Oil Filter Silencer", "/icons/mods/weapon-mod-oilfiltersilencer.avif"},
		{5, "Military Silencer", "/icons/mods/weapon-mod-silencer.avif"},
		{6, "Simple Handmade Sight", "/icons/mods/weapon-mod-simplesight.avif"},
		{7, "Holosight", "/icons/mods/weapon-mod-holosight.avif"},
		{8, "8x Zoom Scope", "/icons/mods/weapon-mod-8x-scope.avif"},
		{9, "Variable Zoom Scope", "/icons/mods/weapon-mod-variable-scope.avif"},
		{10, "Gas Compression Overdrive", "/icons/mods/weapon-mod-gascompressionovedrive.jpeg"},
		{11, "Muzzle Boost", "/icons/mods/weapon-mod-muzzleboost.avif"},
		{12, "Muzzle Brake", "/icons/mods/weapon-mod-muzzlebrake.avif"},
		{13, "Burst Module", "/icons/mods/weapon-mod-burstmodule.avif"},
		{14, "Extended Magazine", "/icons/mods/weapon-mod-extendedmags.avif"},
		{15, "Targeting Attachment", "/icons/mods/weapon-mod-targetingattachment.avif"},
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
		{1, "Metal Fragments", "/icons/ingredients/metal-fragments.avif"},
		{2, "Wood", "/icons/ingredients/wood.avif"},
		{3, "Cloth", "/icons/ingredients/cloth.avif"},
		{4, "High Quality Metal", "/icons/ingredients/metal-refined.avif"},
		{5, "Animal Fat", "/icons/ingredients/leather.avif"},
		{6, "Pipe", "/icons/ingredients/metalpipe.avif"},
		{7, "Spring", "/icons/ingredients/metalspring.avif"},
		{8, "Gear", "/icons/ingredients/gears.avif"},
		{9, "SMG Body", "/icons/ingredients/smgbody.avif"},
		{10, "Semi-Automatic Body", "/icons/ingredients/semibody.avif"},
		{11, "Rifle Body", "/icons/ingredients/riflebody.avif"},
		{12, "Rope", "/icons/ingredients/rope.avif"},
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
			Icon: "/icons/weapons/rifle-ak.avif",
			AmmoIDs: []int{1, 2, 3, 4},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 14, 15},
			Ingredients: []ingredientSeed{{4, 50}, {2, 200}, {7, 4}, {11, 1}},
		},
		{
			ID: 2, Name: "LR-300", Type: "range", Firemode: "automatic", Craftable: false, Stacksize: 1,
			Description: "The LR-300 is a high-accuracy automatic rifle. Cannot be crafted — only found in loot crates.",
			Shortname: "lr300", Capacity: ptr(30), TimeToCraft: nil, CategoryID: 1,
			Icon: "/icons/weapons/rifle-lr300.avif",
			AmmoIDs: []int{1, 2, 3, 4},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 14, 15},
			Ingredients: []ingredientSeed{},
		},
		{
			ID: 3, Name: "MP5A4", Type: "range", Firemode: "automatic", Craftable: true, Stacksize: 1,
			Description: "The MP5A4 is a compact submachine gun with a high rate of fire. Effective in close quarters.",
			Shortname: "mp5a4", Capacity: ptr(30), TimeToCraft: ptr(10), CategoryID: 1,
			Icon: "/icons/weapons/smg-mp5.avif",
			AmmoIDs: []int{5, 6, 7},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 14},
			Ingredients: []ingredientSeed{{4, 15}, {9, 1}, {7, 2}},
		},
		{
			ID: 4, Name: "Thompson", Type: "range", Firemode: "automatic", Craftable: true, Stacksize: 1,
			Description: "The Thompson submachine gun, also known as the \"Tommy Gun\", is a classic automatic weapon.",
			Shortname: "tommy", Capacity: ptr(20), TimeToCraft: ptr(10), CategoryID: 1,
			Icon: "/icons/weapons/smg-thompson.avif",
			AmmoIDs: []int{5, 6, 7},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			Ingredients: []ingredientSeed{{4, 10}, {2, 100}, {9, 1}, {7, 1}},
		},
		{
			ID: 5, Name: "Pump Shotgun", Type: "range", Firemode: "semi", Craftable: true, Stacksize: 1,
			Description: "A pump-action shotgun with a tight spread. Devastating at close range.",
			Shortname: "pump_shotgun", Capacity: ptr(6), TimeToCraft: ptr(15), CategoryID: 1,
			Icon: "/icons/weapons/shotgun-pump.avif",
			AmmoIDs: []int{8, 9, 10, 11},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12},
			Ingredients: []ingredientSeed{{4, 15}, {6, 2}, {7, 1}},
		},
		{
			ID: 6, Name: "Double Barrel Shotgun", Type: "range", Firemode: "double", Craftable: true, Stacksize: 1,
			Description: "A double-barreled shotgun. Two shots, high damage, slow reload.",
			Shortname: "double_barrel", Capacity: ptr(2), TimeToCraft: ptr(10), CategoryID: 1,
			Icon: "/icons/weapons/shotgun-double.avif",
			AmmoIDs: []int{8, 9, 10, 11},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12},
			Ingredients: []ingredientSeed{{1, 175}, {6, 2}},
		},
		{
			ID: 7, Name: "Revolver", Type: "range", Firemode: "semi", Craftable: true, Stacksize: 1,
			Description: "A six-shot revolver. Slower fire rate but packs a punch.",
			Shortname: "revolver", Capacity: ptr(6), TimeToCraft: ptr(5), CategoryID: 1,
			Icon: "/icons/weapons/pistol-revolver.avif",
			AmmoIDs: []int{5, 6, 7},
			ModIDs:  []int{3, 4, 5, 11, 12},
			Ingredients: []ingredientSeed{{6, 1}, {3, 25}, {1, 125}},
		},
		{
			ID: 8, Name: "Semi-Automatic Rifle", Type: "range", Firemode: "semi", Craftable: true, Stacksize: 1,
			Description: "A semi-automatic pistol. Reliable and easy to craft.",
			Shortname: "semi_auto_pistol", Capacity: ptr(8), TimeToCraft: ptr(5), CategoryID: 1,
			Icon: "/icons/weapons/rifle-semiauto.avif",
			AmmoIDs: []int{1, 2, 3, 4},
			ModIDs:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 15},
			Ingredients: []ingredientSeed{{10, 1}, {7, 1}, {1, 450}, {4, 4}},
		},
		{
			ID: 9, Name: "Machete", Type: "melee", Firemode: "none", Craftable: true, Stacksize: 1,
			Description: "A sharp machete. No ammo needed, just swing.",
			Shortname: "machete", Capacity: nil, TimeToCraft: ptr(5), CategoryID: 1,
			Icon: "/icons/weapons/machete.avif",
			AmmoIDs: []int{},
			ModIDs:  []int{},
			Ingredients: []ingredientSeed{{2, 40}, {1, 100}},
		},
		{
			ID: 10, Name: "Hunting Bow", Type: "range", Firemode: "none", Craftable: true, Stacksize: 1,
			Description: "A simple hunting bow. Silent but deadly with practice.",
			Shortname: "hunting_bow", Capacity: ptr(1), TimeToCraft: ptr(5), CategoryID: 1,
			Icon: "/icons/weapons/bow-hunting.avif",
			AmmoIDs: []int{12, 13, 14, 15},
			ModIDs:  []int{},
			Ingredients: []ingredientSeed{{2, 200}, {3, 50}},
		},
		{
			ID: 11, Name: "Compound Bow", Type: "range", Firemode: "none", Craftable: true, Stacksize: 1,
			Description: "A modern compound bow. Higher damage and faster arrow velocity than the standard bow.",
			Shortname: "compound_bow", Capacity: ptr(1), TimeToCraft: ptr(10), CategoryID: 1,
			Icon: "/icons/weapons/bow-compound.avif",
			AmmoIDs: []int{12, 13, 14, 15},
			ModIDs:  []int{},
			Ingredients: []ingredientSeed{{2, 100}, {1, 75}, {12, 2}},
		},
	}

	for _, w := range weapons {
		if _, err := db.Exec(`
			INSERT INTO weapon_item (id, name, type, firemode, craftable, stacksize, description, shortname, icon, capacity, time_to_craft, category_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
			ON CONFLICT (id) DO UPDATE SET name=$2, type=$3, firemode=$4, craftable=$5, stacksize=$6, description=$7, shortname=$8, icon=$9, capacity=$10, time_to_craft=$11, category_id=$12`,
			w.ID, w.Name, w.Type, w.Firemode, w.Craftable, w.Stacksize, w.Description, w.Shortname, w.Icon, w.Capacity, w.TimeToCraft, w.CategoryID,
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

	tables := []string{"category", "ammo", "mods", "ingredients", "weapon_item", "clothing_item", "users", "weapon_suggestions"}
	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf("SELECT setval('%s_id_seq', COALESCE((SELECT MAX(id) FROM %s), 1))", table, table)); err != nil {
			return fmt.Errorf("reset seq %s: %w", table, err)
		}
	}

	log.Info("Database seeded successfully")
	return nil
}

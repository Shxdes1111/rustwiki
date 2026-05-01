package models

type WeaponItem struct {
	ID          int            `json:"id"`
	Name       string        `json:"name"`
	Type       string        `json:"type"`
	Firemode   string        `json:"firemode"`
	Craftable  bool          `json:"craftable"`
	Stacksize  int           `json:"stacksize"`
	CategoryID *int          `json:"category_id,omitempty"`
	Ammo       []Ammo        `json:"ammo,omitempty"`
	Mods       []Mods        `json:"mods,omitempty"`
	Ingredients []Ingredients `json:"ingredients,omitempty"`
}

type Ammo struct {
	ID           int    `json:"id"`
	Name        string `json:"name"`
	WeaponItemID *int  `json:"weapon_item_id,omitempty"`
}

type Mods struct {
	ID           int    `json:"id"`
	Name        string `json:"name"`
	WeaponItemID *int  `json:"weapon_item_id,omitempty"`
}
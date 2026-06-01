package models

type WeaponItem struct {
	ID          int            `json:"id"`
	Name       string        `json:"name"`
	Type       string        `json:"type"`
	Firemode   string        `json:"firemode"`
	Craftable  bool          `json:"craftable"`
	Stacksize  int           `json:"stacksize"`
	Description string       `json:"description,omitempty"`
	Shortname  string        `json:"shortname,omitempty"`
	Capacity   *int          `json:"capacity,omitempty"`
	TimeToCraft *int         `json:"time_to_craft,omitempty"`
	CategoryID *int          `json:"category_id,omitempty"`
	Ammo       []Ammo        `json:"ammo,omitempty"`
	Mods       []Mods        `json:"mods,omitempty"`
	Ingredients []Ingredients `json:"ingredients,omitempty"`
}

type Ammo struct {
	ID           int    `json:"id"`
	Name        string `json:"name"`
	Icon        *string `json:"icon,omitempty"`
	WeaponItemID *int  `json:"weapon_item_id,omitempty"`
}

type Mods struct {
	ID           int    `json:"id"`
	Name        string `json:"name"`
	Icon        *string `json:"icon,omitempty"`
	WeaponItemID *int  `json:"weapon_item_id,omitempty"`
}
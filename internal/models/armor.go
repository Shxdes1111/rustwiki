package models

type ClothingItem struct {
	ID             int            `json:"id"`
	Name           string        `json:"name"`
	Equipmentslot  string        `json:"equipmentslot"`
	Protection     *int          `json:"protection,omitempty"`
	Craftable      bool          `json:"craftable"`
	Stacksize      int           `json:"stacksize"`
	CategoryID    *int          `json:"category_id,omitempty"`
	Ingredients   []Ingredients `json:"ingredients,omitempty"`
}

type Ingredients struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	WeaponItemID   *int   `json:"weapon_item_id,omitempty"`
	ClothingItemID *int   `json:"clothing_item_id,omitempty"`
}
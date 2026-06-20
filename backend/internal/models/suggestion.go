package models

import (
	"encoding/json"
	"time"
)

type WeaponSuggestion struct {
	ID         int              `json:"id"`
	UserID     int              `json:"user_id"`
	Username   string           `json:"username,omitempty"`
	Payload    json.RawMessage  `json:"payload"`
	Status     string           `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	ReviewedAt *time.Time       `json:"reviewed_at,omitempty"`
	ReviewedBy *int             `json:"reviewed_by,omitempty"`
}

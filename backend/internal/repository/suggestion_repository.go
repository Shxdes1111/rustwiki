package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"backend/internal/logger"
	"backend/internal/models"
)

type SuggestionRepository interface {
	Create(userID int, payload json.RawMessage) (*models.WeaponSuggestion, error)
	FindAll() ([]models.WeaponSuggestion, error)
	FindByID(id int) (*models.WeaponSuggestion, error)
	UpdateStatus(id, reviewedBy int, status string) error
	RemoveIconBase64(id int) error
	Delete(id int) error
}

type suggestionRepository struct {
	db  *sql.DB
	log *logger.Logger
}

func NewSuggestionRepository(db *sql.DB, log *logger.Logger) SuggestionRepository {
	return &suggestionRepository{db: db, log: log}
}

func (r *suggestionRepository) Create(userID int, payload json.RawMessage) (*models.WeaponSuggestion, error) {
	var s models.WeaponSuggestion
	err := r.db.QueryRow(
		`INSERT INTO weapon_suggestions (user_id, payload) VALUES ($1, $2)
		 RETURNING id, user_id, payload, status, created_at`,
		userID, payload,
	).Scan(&s.ID, &s.UserID, &s.Payload, &s.Status, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *suggestionRepository) FindAll() ([]models.WeaponSuggestion, error) {
	rows, err := r.db.Query(`
		SELECT s.id, s.user_id, COALESCE(u.username, ''), s.payload, s.status, s.created_at, s.reviewed_at, s.reviewed_by
		FROM weapon_suggestions s
		LEFT JOIN users u ON s.user_id = u.id
		ORDER BY s.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []models.WeaponSuggestion
	for rows.Next() {
		var s models.WeaponSuggestion
		if err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Payload, &s.Status, &s.CreatedAt, &s.ReviewedAt, &s.ReviewedBy); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, s)
	}
	return suggestions, nil
}

func (r *suggestionRepository) FindByID(id int) (*models.WeaponSuggestion, error) {
	var s models.WeaponSuggestion
	err := r.db.QueryRow(`
		SELECT s.id, s.user_id, COALESCE(u.username, ''), s.payload, s.status, s.created_at, s.reviewed_at, s.reviewed_by
		FROM weapon_suggestions s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.id = $1
	`, id).Scan(&s.ID, &s.UserID, &s.Username, &s.Payload, &s.Status, &s.CreatedAt, &s.ReviewedAt, &s.ReviewedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *suggestionRepository) UpdateStatus(id, reviewedBy int, status string) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE weapon_suggestions SET status = $1, reviewed_at = $2, reviewed_by = $3 WHERE id = $4`,
		status, now, reviewedBy, id,
	)
	return err
}

func (r *suggestionRepository) Delete(id int) error {
	_, err := r.db.Exec(
		`DELETE FROM weapon_suggestions WHERE id = $1 AND status != 'pending'`,
		id,
	)
	return err
}

func (r *suggestionRepository) RemoveIconBase64(id int) error {
	_, err := r.db.Exec(
		`UPDATE weapon_suggestions SET payload = payload - 'icon_base64' WHERE id = $1`,
		id,
	)
	return err
}

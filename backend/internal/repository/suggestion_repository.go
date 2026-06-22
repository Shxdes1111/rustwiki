package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/logger"
	"backend/internal/models"
)

type SuggestionRepository interface {
	Create(userID int, payload json.RawMessage) (*models.WeaponSuggestion, error)
	FindAll() ([]models.WeaponSuggestion, error)
	FindByID(id int) (*models.WeaponSuggestion, error)
	FindByUserID(userID int) ([]models.WeaponSuggestion, error)
	UpdateStatus(id, reviewedBy int, status string, rejectionReason *string) error
	RemoveIconBase64(id int) error
	UpdatePayload(id int, payload json.RawMessage) error
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

var suggestionColumns = `s.id, s.user_id, COALESCE(u.username, ''), s.payload, s.status, s.created_at, s.reviewed_at, s.reviewed_by, s.rejection_reason`

func scanSuggestion(scanner interface {
	Scan(dest ...interface{}) error
}, s *models.WeaponSuggestion) error {
	return scanner.Scan(&s.ID, &s.UserID, &s.Username, &s.Payload, &s.Status, &s.CreatedAt, &s.ReviewedAt, &s.ReviewedBy, &s.RejectionReason)
}

func (r *suggestionRepository) FindAll() ([]models.WeaponSuggestion, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM weapon_suggestions s
		LEFT JOIN users u ON s.user_id = u.id
		ORDER BY s.created_at DESC
	`, suggestionColumns)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []models.WeaponSuggestion
	for rows.Next() {
		var s models.WeaponSuggestion
		if err := scanSuggestion(rows, &s); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return suggestions, nil
}

func (r *suggestionRepository) FindByID(id int) (*models.WeaponSuggestion, error) {
	var s models.WeaponSuggestion
	query := fmt.Sprintf(`
		SELECT %s
		FROM weapon_suggestions s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.id = $1
	`, suggestionColumns)
	err := scanSuggestion(r.db.QueryRow(query, id), &s)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *suggestionRepository) FindByUserID(userID int) ([]models.WeaponSuggestion, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM weapon_suggestions s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.user_id = $1
		ORDER BY s.created_at DESC
	`, suggestionColumns)
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []models.WeaponSuggestion
	for rows.Next() {
		var s models.WeaponSuggestion
		if err := scanSuggestion(rows, &s); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return suggestions, nil
}

func (r *suggestionRepository) UpdateStatus(id, reviewedBy int, status string, rejectionReason *string) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE weapon_suggestions SET status = $1, reviewed_at = $2, reviewed_by = $3, rejection_reason = $4 WHERE id = $5`,
		status, now, reviewedBy, rejectionReason, id,
	)
	return err
}

func (r *suggestionRepository) UpdatePayload(id int, payload json.RawMessage) error {
	_, err := r.db.Exec(
		`UPDATE weapon_suggestions SET payload = $1, status = 'pending', reviewed_at = NULL, reviewed_by = NULL, rejection_reason = NULL WHERE id = $2`,
		payload, id,
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

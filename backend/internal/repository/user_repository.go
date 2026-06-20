package repository

import (
	"database/sql"

	"backend/internal/logger"
	"backend/internal/models"
)

type UserRepository interface {
	Create(username, passwordHash, role string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByID(id int) (*models.User, error)
}

type userRepository struct {
	db  *sql.DB
	log *logger.Logger
}

func NewUserRepository(db *sql.DB, log *logger.Logger) UserRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) Create(username, passwordHash, role string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3)
		 RETURNING id, username, role, created_at`,
		username, passwordHash, role,
	).Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, role, created_at FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, role, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

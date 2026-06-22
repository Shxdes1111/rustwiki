package database

import (
	"database/sql"
	"fmt"
	"time"

	"backend/internal/config"
	"backend/internal/logger"
)

func AutoMigrate(db *sql.DB, log *logger.Logger) error {
	log.Info("Running auto-migration...")

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('guest', 'user', 'admin')),
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS weapon_suggestions (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			payload JSONB NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
			created_at TIMESTAMP DEFAULT NOW(),
			reviewed_at TIMESTAMP,
			reviewed_by INTEGER REFERENCES users(id),
			rejection_reason TEXT
		)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("auto-migrate: %w", err)
		}
	}

	if _, err := db.Exec(`ALTER TABLE weapon_suggestions ADD COLUMN IF NOT EXISTS rejection_reason TEXT`); err != nil {
		return fmt.Errorf("auto-migrate add rejection_reason: %w", err)
	}

	if _, err := db.Exec(`ALTER TABLE weapon_item ADD COLUMN IF NOT EXISTS views INTEGER NOT NULL DEFAULT 0`); err != nil {
		return fmt.Errorf("auto-migrate add views: %w", err)
	}
	if _, err := db.Exec(`ALTER TABLE weapon_item ADD COLUMN IF NOT EXISTS created_by INTEGER REFERENCES users(id)`); err != nil {
		return fmt.Errorf("auto-migrate add created_by: %w", err)
	}

	log.Info("Auto-migration completed")
	return nil
}

// DB представляет подключение к базе данных
type DB struct {
	*sql.DB
}

// Connect создает подключение к базе данных
func Connect(cfg *config.DatabaseConfig, log *logger.Logger) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)                 // Максимальное количество открытых соединений
	db.SetMaxIdleConns(5)                  // Максимальное количество неактивных соединений
	db.SetConnMaxLifetime(5 * time.Minute) // Максимальное время жизни соединения

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Successfully connected to database")

	return &DB{DB: db}, nil
}

// Close закрывает подключение к базе данных
func (db *DB) Close() error {
	return db.DB.Close()
}

// Health проверяет состояние базы данных
func (db *DB) Health() error {
	return db.Ping()
}

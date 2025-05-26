package database

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     int
	SSLMode  string
}

func NewSQLDatabase(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString(cfg))
	if err != nil {
		return nil, fmt.Errorf("sql.Open() returned error: %w", err)
	}

	return db, nil
}

func connString(cfg Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
}

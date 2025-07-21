package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DB wraps sql.DB to provide additional functionality
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection and verifies connectivity
func NewDB(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

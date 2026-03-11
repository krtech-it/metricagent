package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewDB(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("dsn is empty")
	}

	if err := applyMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func applyMigrations(dsn string) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()
	return m.Up()
}

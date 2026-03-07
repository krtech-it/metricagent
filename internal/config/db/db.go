package db

import (
	"database/sql"
	"errors"
)

func NewDB(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("dsn is empty")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	InitSchema(db)
	return db, nil
}

func InitSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE metrics (
    id VARCHAR(50) PRIMARY KEY,
    m_type VARCHAR(50) NOT NULL,
    delta bigint default null,
    value double precision default null
)
	`)
	return err
}

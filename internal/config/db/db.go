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
	return db, nil
}

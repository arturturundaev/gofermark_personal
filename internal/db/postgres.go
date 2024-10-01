package db

import (
	"github.com/jmoiron/sqlx"
	"time"
)

func NewPostgres(dsn string, maxOpen int, maxIdle int, maxLifeDuration time.Duration) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(maxLifeDuration)

	return db, nil
}

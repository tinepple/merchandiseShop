package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*Storage, error) {
	return &Storage{
		db: db,
	}, nil
}

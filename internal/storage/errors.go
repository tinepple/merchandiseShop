package storage

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("row not found")

func handleSQLError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

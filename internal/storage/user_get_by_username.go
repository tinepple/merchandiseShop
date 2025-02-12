package storage

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

var ErrNotFound = errors.New("row not found")

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (User, error) {
	query, params, err := sq.Select(
		"id",
		"username",
		"password",
	).From(usersTableName).
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return User{}, err
	}

	var dest User

	err = s.db.QueryRowContext(ctx, s.db.Rebind(query), params...).Scan(
		&dest.ID,
		&dest.Username,
		&dest.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return dest, nil
}

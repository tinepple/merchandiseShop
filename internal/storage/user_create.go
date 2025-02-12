package storage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) CreateUser(ctx context.Context, username string, password string) (User, error) {
	query, params, err := sq.Insert(usersTableName).
		Columns(
			"username",
			"password",
		).
		Values(username, password).
		Suffix("returning id, username,password").
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
		return User{}, err
	}

	return dest, nil
}

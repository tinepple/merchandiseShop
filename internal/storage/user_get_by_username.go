package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
)

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
		return User{}, handleSQLError(err)
	}

	return dest, nil
}

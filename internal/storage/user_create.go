package storage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) CreateUser(ctx context.Context, username string, password string) (User, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return User{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
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

	err = s.db.QueryRowContext(ctx, tx.Rebind(query), params...).Scan(
		&dest.ID,
		&dest.Username,
		&dest.Password,
	)
	if err != nil {
		return User{}, err
	}

	query, params, err = sq.Insert(balancesTableName).
		Columns(
			"user_id",
		).
		Values(dest.ID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return User{}, err
	}

	_, err = tx.ExecContext(ctx, tx.Rebind(query), params...)
	if err != nil {
		return User{}, err
	}

	err = tx.Commit()
	if err != nil {
		return User{}, err
	}

	return dest, nil
}

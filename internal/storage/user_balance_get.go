package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetUserBalance(ctx context.Context, userID int) (int, error) {
	query, params, err := sq.Select("balance").
		From(balancesTableName).
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	var balance int
	err = s.db.QueryRowContext(ctx, s.db.Rebind(query), params...).Scan(&balance)
	if err != nil {
		return 0, handleSQLError(err)
	}

	return balance, nil
}

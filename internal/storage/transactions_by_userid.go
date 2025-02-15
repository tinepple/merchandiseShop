package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetTransactionsByUserID(ctx context.Context, userID int) ([]CoinsHistory, error) {
	query, params, err := sq.Select(
		"t.user_id_from as user_id_from",
		"uf.username as username_from",
		"t.user_id_to user_id_to",
		"ut.username username_to",
		"t.amount as amount",
	).From(fmt.Sprintf("%s t", transactionsTableName)).
		InnerJoin(fmt.Sprintf("%s uf on uf.id = t.user_id_from", usersTableName)).
		InnerJoin(fmt.Sprintf("%s ut on ut.id = t.user_id_to", usersTableName)).
		Where(sq.Or{
			sq.Eq{"t.user_id_from": userID},
			sq.Eq{"t.user_id_to": userID},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var dest []CoinsHistory

	err = s.db.SelectContext(ctx, &dest, s.db.Rebind(query), params...)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

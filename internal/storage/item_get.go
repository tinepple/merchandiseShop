package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetItem(ctx context.Context, name string) (Item, error) {
	query, params, err := sq.Select("id, name, price").
		From(merchandiseTableName).
		Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return Item{}, err
	}

	var item Item

	err = s.db.QueryRowContext(ctx, s.db.Rebind(query), params...).Scan(
		&item.ID,
		&item.Name,
		&item.Price,
	)
	if err != nil {
		return Item{}, handleSQLError(err)
	}

	return item, nil
}

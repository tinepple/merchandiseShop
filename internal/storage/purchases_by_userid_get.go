package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetPurchasesByUserID(ctx context.Context, userID int) ([]Inventory, error) {
	query, params, err := sq.Select(
		"m.name as name",
		"count(*) as quantity",
	).From(fmt.Sprintf("%s p", purchasesTableName)).
		InnerJoin(fmt.Sprintf("%s m on m.id = p.merchandise_id", merchandiseTableName)).
		GroupBy("m.name").
		Where(sq.Eq{"p.user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var dest []Inventory

	err = s.db.SelectContext(ctx, &dest, s.db.Rebind(query), params...)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

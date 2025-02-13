package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) CreatePurchase(ctx context.Context, userID, itemID, newBalance int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = s.updateUserBalance(ctx, tx, userID, newBalance)
	if err != nil {
		return err
	}

	err = s.createPurchase(ctx, tx, userID, itemID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) updateUserBalance(ctx context.Context, tx *sqlx.Tx, userID int, newBalance int) error {
	query, params, err := sq.Update(balancesTableName).
		Set(
			"balance", newBalance).
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, tx.Rebind(query), params...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) createPurchase(ctx context.Context, tx *sqlx.Tx, userID, itemID int) error {
	query, params, err := sq.Insert(purchasesTableName).
		Columns(
			"user_id",
			"merchandise_id",
		).
		Values(userID, itemID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, tx.Rebind(query), params...)
	if err != nil {
		return err
	}

	return nil
}

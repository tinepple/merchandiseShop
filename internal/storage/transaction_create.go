package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) CreateTransaction(ctx context.Context, transaction Transaction) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = s.createTransaction(ctx, tx, transaction)
	if err != nil {
		return err
	}

	err = s.updateUserBalance(ctx, tx, transaction.UserIDFrom, transaction.NewBalanceUserFrom)
	if err != nil {
		return err
	}

	err = s.updateUserBalance(ctx, tx, transaction.UserIDTo, transaction.NewBalanceUserTo)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) createTransaction(ctx context.Context, tx *sqlx.Tx, transaction Transaction) error {
	query, params, err := sq.Insert(transactionsTableName).
		Columns(
			"user_id_from",
			"user_id_to",
			"amount",
		).
		Values(
			transaction.UserIDFrom,
			transaction.UserIDTo,
			transaction.Amount,
		).
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

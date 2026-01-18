package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	DB *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) Save(
	ctx context.Context,
	tx *sqlx.Tx,
	fromAccountID string,
	toAccountID string,
	amount int64,
	txType string,
) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO transactions
		(from_account_id, to_account_id, amount, type)
		VALUES ($1, $2, $3, $4)
	`,
		fromAccountID,
		toAccountID,
		amount,
		txType,
	)
	return err
}

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
	fromAccount,
	toAccount string,
	amount int64,
	txType string,
) error {

	query := `
		INSERT INTO transactions (from_account, to_account, amount, transaction_type)
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		fromAccount,
		toAccount,
		amount,
		txType,
	)

	return err
}

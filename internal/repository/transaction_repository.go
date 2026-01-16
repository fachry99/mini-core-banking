package repository

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type TransferRepository struct {
	DB *sqlx.DB
}

func NewTransferRepository(db *sqlx.DB) *TransferRepository {
	return &TransferRepository{DB: db}
}

func (r *TransferRepository) Transfer(
	ctx context.Context,
	fromAccountID string,
	toAccountID string,
	amount int64,
) error {

	if amount <= 0 {
		return errors.New("invalid transfer amount")
	}

	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Safety rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var fromBalance int64

	// 🔒 Lock sender account
	err = tx.GetContext(ctx, &fromBalance, `
		SELECT balance FROM accounts
		WHERE id = $1
		FOR UPDATE
	`, fromAccountID)
	if err != nil {
		return err
	}

	if fromBalance < amount {
		return errors.New("insufficient balance")
	}

	// Deduct sender
	_, err = tx.ExecContext(ctx, `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
	`, amount, fromAccountID)
	if err != nil {
		return err
	}

	// Add receiver
	_, err = tx.ExecContext(ctx, `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
	`, amount, toAccountID)
	if err != nil {
		return err
	}

	// Save transaction history
	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions
		(from_account_id, to_account_id, amount, type)
		VALUES ($1, $2, $3, 'TRANSFER')
	`,
		fromAccountID,
		toAccountID,
		amount,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

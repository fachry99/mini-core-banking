// TODO: will be used in phase service refactor

package service

import (
	"context"
	"errors"

	"github.com/fachry/mini-core-banking/internal/repository"
	"github.com/jmoiron/sqlx"
)

type TransferService struct {
	DB              *sqlx.DB
	AccountRepo     *repository.AccountRepository
	TransactionRepo *repository.TransactionRepository
}

func NewTransferService(
	db *sqlx.DB,
	accountRepo *repository.AccountRepository,
	transactionRepo *repository.TransactionRepository,
) *TransferService {
	return &TransferService{
		DB:              db,
		AccountRepo:     accountRepo,
		TransactionRepo: transactionRepo,
	}
}
func (s *TransferService) Transfer(
	ctx context.Context,
	fromAccountID string,
	toAccountID string,
	amount int64,
) error {

	if amount <= 0 {
		return errors.New("invalid transfer amount")
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Safety net
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
		UPDATE accounts SET balance = balance - $1
		WHERE id = $2
	`, amount, fromAccountID)
	if err != nil {
		return err
	}

	// Add receiver
	_, err = tx.ExecContext(ctx, `
		UPDATE accounts SET balance = balance + $1
		WHERE id = $2
	`, amount, toAccountID)
	if err != nil {
		return err
	}

	// Save transaction record
	err = s.TransactionRepo.Save(
		ctx,
		tx,
		fromAccountID,
		toAccountID,
		amount,
		"TRANSFER",
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

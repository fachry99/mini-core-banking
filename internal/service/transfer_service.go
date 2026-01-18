package service

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/fachry/mini-core-banking/internal/audit"
	"github.com/fachry/mini-core-banking/internal/middleware"
)

// =====================
// DOMAIN ERRORS
// =====================
var (
	ErrSameAccount         = errors.New("cannot transfer to same account")
	ErrInvalidAmount       = errors.New("invalid transfer amount")
	ErrAccountNotFound     = errors.New("account not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type TransferService struct {
	DB *sqlx.DB
}

func NewTransferService(db *sqlx.DB) *TransferService {
	return &TransferService{DB: db}
}

func (s *TransferService) Transfer(
	ctx context.Context,
	fromAccountID string,
	toAccountID string,
	amount int64,
) (err error) {

	requestID, _ := ctx.Value(middleware.RequestIDKey).(string)

	// =====================
	// 🔍 AUDIT FINALIZER
	// =====================
	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}

		audit.LogTransfer(
			requestID,
			fromAccountID,
			toAccountID,
			amount,
			status,
		)
	}()

	// =====================
	// 1️⃣ BASIC VALIDATION
	// =====================
	if fromAccountID == toAccountID {
		return ErrSameAccount
	}

	if amount <= 0 {
		return ErrInvalidAmount
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// =====================
	// 2️⃣ CONSISTENT LOCK ORDER
	// =====================
	firstID, secondID := fromAccountID, toAccountID
	if firstID > secondID {
		firstID, secondID = secondID, firstID
	}

	type account struct {
		ID      string `db:"id"`
		Balance int64  `db:"balance"`
	}

	accounts := map[string]*account{}

	rows, err := tx.QueryxContext(ctx, `
		SELECT id, balance
		FROM accounts
		WHERE id IN ($1, $2)
		FOR UPDATE
	`, firstID, secondID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var acc account
		if err := rows.StructScan(&acc); err != nil {
			return err
		}
		accounts[acc.ID] = &acc
	}

	if len(accounts) != 2 {
		return ErrAccountNotFound
	}

	// =====================
	// 3️⃣ BUSINESS RULE
	// =====================
	if accounts[fromAccountID].Balance < amount {
		return ErrInsufficientBalance
	}

	// =====================
	// 4️⃣ UPDATE BALANCES
	// =====================
	if _, err = tx.ExecContext(ctx, `
		UPDATE accounts SET balance = balance - $1 WHERE id = $2
	`, amount, fromAccountID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE accounts SET balance = balance + $1 WHERE id = $2
	`, amount, toAccountID); err != nil {
		return err
	}

	// =====================
	// 5️⃣ TRANSACTION LOG
	// =====================
	if _, err = tx.ExecContext(ctx, `
		INSERT INTO transactions
		(from_account_id, to_account_id, amount, type)
		VALUES ($1, $2, $3, 'TRANSFER')
	`, fromAccountID, toAccountID, amount); err != nil {
		return err
	}

	// =====================
	// 6️⃣ COMMIT
	// =====================
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

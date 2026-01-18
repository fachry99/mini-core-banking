// package service

// import (
// 	"context"
// 	"errors"

// 	"github.com/jmoiron/sqlx"
// )

// type TransferService struct {
// 	DB *sqlx.DB
// }

// func NewTransferService(db *sqlx.DB) *TransferService {
// 	return &TransferService{DB: db}
// }

// func (s *TransferService) Transfer(
// 	ctx context.Context,
// 	fromAccountID string,
// 	toAccountID string,
// 	amount int64,
// ) error {

// 	// =====================
// 	// 1️⃣ BASIC VALIDATION
// 	// =====================
// 	if fromAccountID == toAccountID {
// 		return errors.New("cannot transfer to same account")
// 	}

// 	if amount <= 0 {
// 		return errors.New("invalid transfer amount")
// 	}

// 	tx, err := s.DB.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// =====================
// 	// 2️⃣ CONSISTENT LOCK ORDER (ANTI DEADLOCK)
// 	// =====================
// 	firstID := fromAccountID
// 	secondID := toAccountID
// 	if firstID > secondID {
// 		firstID, secondID = secondID, firstID
// 	}

// 	type account struct {
// 		ID      string `db:"id"`
// 		Balance int64  `db:"balance"`
// 	}

// 	accounts := map[string]*account{}

// 	rows, err := tx.QueryxContext(ctx, `
// 		SELECT id, balance
// 		FROM accounts
// 		WHERE id IN ($1, $2)
// 		FOR UPDATE
// 	`, firstID, secondID)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var acc account
// 		if err := rows.StructScan(&acc); err != nil {
// 			return err
// 		}
// 		accounts[acc.ID] = &acc
// 	}

// 	if len(accounts) != 2 {
// 		return errors.New("account not found")
// 	}

// 	// =====================
// 	// 3️⃣ BUSINESS RULE
// 	// =====================
// 	if accounts[fromAccountID].Balance < amount {
// 		return errors.New("insufficient balance")
// 	}

// 	// =====================
// 	// 4️⃣ UPDATE BALANCES
// 	// =====================
// 	_, err = tx.ExecContext(ctx, `
// 		UPDATE accounts
// 		SET balance = balance - $1
// 		WHERE id = $2
// 	`, amount, fromAccountID)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = tx.ExecContext(ctx, `
// 		UPDATE accounts
// 		SET balance = balance + $1
// 		WHERE id = $2
// 	`, amount, toAccountID)
// 	if err != nil {
// 		return err
// 	}

// 	// =====================
// 	// 5️⃣ TRANSACTION LOG
// 	// =====================
// 	_, err = tx.ExecContext(ctx, `
// 		INSERT INTO transactions
// 		(from_account_id, to_account_id, amount, type)
// 		VALUES ($1, $2, $3, 'TRANSFER')
// 	`, fromAccountID, toAccountID, amount)
// 	if err != nil {
// 		return err
// 	}

//		// =====================
//		// 6️⃣ COMMIT
//		// =====================
//		return tx.Commit()
//	}
package service

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
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
) error {

	// =====================
	// 1️⃣ BASIC VALIDATION
	// =====================
	if fromAccountID == toAccountID {
		return errors.New("cannot transfer to same account")
	}

	if amount <= 0 {
		return errors.New("invalid transfer amount")
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// =====================
	// 2️⃣ CONSISTENT LOCK ORDER (ANTI DEADLOCK)
	// =====================
	firstID := fromAccountID
	secondID := toAccountID
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
		return errors.New("account not found")
	}

	// =====================
	// 3️⃣ BUSINESS RULE
	// =====================
	if accounts[fromAccountID].Balance < amount {
		return errors.New("insufficient balance")
	}

	// =====================
	// 4️⃣ UPDATE BALANCES
	// =====================
	_, err = tx.ExecContext(ctx, `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
	`, amount, fromAccountID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
	`, amount, toAccountID)
	if err != nil {
		return err
	}

	// =====================
	// 5️⃣ TRANSACTION LOG
	// =====================
	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions
		(from_account_id, to_account_id, amount, type)
		VALUES ($1, $2, $3, 'TRANSFER')
	`, fromAccountID, toAccountID, amount)
	if err != nil {
		return err
	}

	// =====================
	// 6️⃣ COMMIT
	// =====================
	return tx.Commit()
}

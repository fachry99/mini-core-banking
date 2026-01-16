package repository

import (
	"context"
	"errors"

	"github.com/fachry/mini-core-banking/internal/domain"

	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	DB *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `
		INSERT INTO accounts (user_id, account_number, balance)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return r.DB.QueryRowContext(
		ctx,
		query,
		account.UserID,
		account.AccountNumber,
		account.Balance,
	).Scan(&account.ID, &account.CreatedAt)
}

func (r *AccountRepository) Deposit(
	ctx context.Context,
	accountID string,
	amount int64,
) error {

	if amount <= 0 {
		return errors.New("invalid deposit amount")
	}

	_, err := r.DB.ExecContext(ctx, `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
	`, amount, accountID)

	return err
}

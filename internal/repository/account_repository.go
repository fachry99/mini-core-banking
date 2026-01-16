package repository

import (
	"context"

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

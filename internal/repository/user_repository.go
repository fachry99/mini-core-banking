package repository

import (
	"context"

	"github.com/fachry/mini-core-banking/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (full_name, email)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	return r.DB.QueryRowContext(
		ctx,
		query,
		user.FullName,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)
}

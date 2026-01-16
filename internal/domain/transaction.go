package domain

import "time"

type Transaction struct {
	ID            string    `db:"id" json:"id"`
	FromAccountID *string   `db:"from_account_id" json:"from_account_id"`
	ToAccountID   *string   `db:"to_account_id" json:"to_account_id"`
	Amount        int64     `db:"amount" json:"amount"`
	Type          string    `db:"type" json:"type"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

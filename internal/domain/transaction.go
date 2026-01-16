package domain

import "time"

type Transaction struct {
	ID              string    `db:"id"`
	FromAccount     *string   `db:"from_account"`
	ToAccount       *string   `db:"to_account"`
	Amount          int64     `db:"amount"`
	TransactionType string    `db:"transaction_type"`
	CreatedAt       time.Time `db:"created_at"`
}

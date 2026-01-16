package domain

import "time"

type Account struct {
	ID            string    `db:"id"`
	UserID        string    `db:"user_id"`
	AccountNumber string    `db:"account_number"`
	Balance       int64     `db:"balance"`
	CreatedAt     time.Time `db:"created_at"`
}

package dto

import "time"

type CreateAccountRequest struct {
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

type AccountResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

package dto

type CreateAccountRequest struct {
	UserID        string `json:"user_id"`
	AccountNumber string `json:"account_number"`
}

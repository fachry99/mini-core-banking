package dto

import "errors"

type CreateAccountRequest struct {
	UserID        string `json:"user_id"`
	AccountNumber string `json:"account_number"`
}

func (r *CreateAccountRequest) Validate() error {
	if r.UserID == "" || r.AccountNumber == "" {
		return errors.New("user_id and account_number are required")
	}
	return nil
}

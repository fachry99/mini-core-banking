package dto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type TransferRequest struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
}

func (r *TransferRequest) Validate() error {
	if r.FromAccountID == "" || r.ToAccountID == "" {
		return errors.New("account id is required")
	}

	if r.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if r.FromAccountID == r.ToAccountID {
		return errors.New("cannot transfer to same account")
	}

	return nil
}

func (r TransferRequest) Hash() string {
	b, _ := json.Marshal(r)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

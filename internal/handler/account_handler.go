package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fachry/mini-core-banking/internal/domain"
	"github.com/fachry/mini-core-banking/internal/repository"
)

type AccountHandler struct {
	Repo *repository.AccountRepository
}

func NewAccountHandler(repo *repository.AccountRepository) *AccountHandler {
	return &AccountHandler{Repo: repo}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID        string `json:"user_id"`
		AccountNumber string `json:"account_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	account := &domain.Account{
		UserID:        req.UserID,
		AccountNumber: req.AccountNumber,
		Balance:       0,
	}

	if err := h.Repo.Create(ctx, account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

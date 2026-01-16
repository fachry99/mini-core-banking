package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fachry/mini-core-banking/internal/repository"
)

type DepositHandler struct {
	Repo *repository.AccountRepository
}

func NewDepositHandler(repo *repository.AccountRepository) *DepositHandler {
	return &DepositHandler{Repo: repo}
}

func (h *DepositHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID string `json:"account_id"`
		Amount    int64  `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.Repo.Deposit(ctx, req.AccountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "deposit success",
	})
}

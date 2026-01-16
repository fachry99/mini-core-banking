package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fachry/mini-core-banking/internal/repository"
)

type TransferHandler struct {
	Repo *repository.TransferRepository
}

func NewTransferHandler(repo *repository.TransferRepository) *TransferHandler {
	return &TransferHandler{Repo: repo}
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromAccountID string `json:"from_account_id"`
		ToAccountID   string `json:"to_account_id"`
		Amount        int64  `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "amount must be positive", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.Repo.Transfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "transfer success",
	})
}

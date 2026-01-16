package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fachry/mini-core-banking/internal/service"
)

type TransferHandler struct {
	Service *service.TransferService
}

func NewTransferHandler(service *service.TransferService) *TransferHandler {
	return &TransferHandler{Service: service}
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

	if err := h.Service.Transfer(
		r.Context(),
		req.FromAccountID,
		req.ToAccountID,
		req.Amount,
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

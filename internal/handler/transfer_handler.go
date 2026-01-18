package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fachry/mini-core-banking/internal/dto"
	"github.com/fachry/mini-core-banking/internal/repository"
	"github.com/fachry/mini-core-banking/internal/service"
)

//	type TransferHandler struct {
//		Service *service.TransferService
//	}
type TransferHandler struct {
	Service         *service.TransferService
	IdempotencyRepo *repository.IdempotencyRepository
}

func NewTransferHandler(
	service *service.TransferService,
	idemRepo *repository.IdempotencyRepository,
) *TransferHandler {
	return &TransferHandler{
		Service:         service,
		IdempotencyRepo: idemRepo,
	}
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// =========================
	// 1️⃣ IDEMPOTENCY KEY (EARLY)
	// =========================
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		http.Error(w, "missing Idempotency-Key", http.StatusBadRequest)
		return
	}

	// kalau key sudah pernah dipakai → RETURN RESPONSE LAMA
	if resp, err := h.IdempotencyRepo.Get(key); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	// =========================
	// 2️⃣ PARSE & VALIDATE BODY
	// =========================
	var req dto.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// =========================
	// 3️⃣ EXECUTE TRANSFER
	// =========================
	if err := h.Service.Transfer(
		ctx,
		req.FromAccountID,
		req.ToAccountID,
		req.Amount,
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// =========================
	// 4️⃣ SAVE IDEMPOTENT RESULT
	// =========================
	response := map[string]string{
		"status": "transfer success",
	}

	responseBytes, _ := json.Marshal(response)

	_ = h.IdempotencyRepo.Save(
		key,
		req.Hash(),
		responseBytes,
	)

	// =========================
	// 5️⃣ RETURN RESPONSE
	// =========================
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

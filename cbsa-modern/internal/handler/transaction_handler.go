package handler

import (
	"net/http"

	"github.com/cicsdev/cbsa-modern/internal/model"
	"github.com/cicsdev/cbsa-modern/internal/service"
)

type TransactionHandler struct {
	transactions *service.TransactionService
}

func NewTransactionHandler(transactions *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactions: transactions}
}

func (h *TransactionHandler) Debit(w http.ResponseWriter, r *http.Request) {
	var req model.DebitCreditRequest
	if err := DecodeJSON(r, &req); err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.AccountNumber == "" || req.Amount == "" {
		ErrorJSON(w, http.StatusBadRequest, "accountNumber and amount are required")
		return
	}

	account, err := h.transactions.Debit(r.Context(), req)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, http.StatusOK, account)
}

func (h *TransactionHandler) Credit(w http.ResponseWriter, r *http.Request) {
	var req model.DebitCreditRequest
	if err := DecodeJSON(r, &req); err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.AccountNumber == "" || req.Amount == "" {
		ErrorJSON(w, http.StatusBadRequest, "accountNumber and amount are required")
		return
	}

	account, err := h.transactions.Credit(r.Context(), req)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, http.StatusOK, account)
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req model.TransferRequest
	if err := DecodeJSON(r, &req); err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.FromAccountNumber == "" || req.ToAccountNumber == "" || req.Amount == "" {
		ErrorJSON(w, http.StatusBadRequest, "fromAccountNumber, toAccountNumber, and amount are required")
		return
	}

	if err := h.transactions.Transfer(r.Context(), req); err != nil {
		ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"message": "transfer successful"})
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	params := model.ListParams{
		Limit:  QueryInt(r, "limit", 20),
		Offset: QueryInt(r, "offset", 0),
	}

	txns, total, err := h.transactions.List(r.Context(), params)
	if err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{
		"transactions": txns,
		"total":        total,
		"limit":        params.Limit,
		"offset":       params.Offset,
	})
}

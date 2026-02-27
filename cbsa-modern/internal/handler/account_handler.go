package handler

import (
	"net/http"

	"github.com/cicsdev/cbsa-modern/internal/model"
	"github.com/cicsdev/cbsa-modern/internal/service"
	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	accounts *service.AccountService
}

func NewAccountHandler(accounts *service.AccountService) *AccountHandler {
	return &AccountHandler{accounts: accounts}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateAccountRequest
	if err := DecodeJSON(r, &req); err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CustomerNumber == "" || req.AccountType == "" {
		ErrorJSON(w, http.StatusBadRequest, "customerNumber and accountType are required")
		return
	}

	account, err := h.accounts.Create(r.Context(), req)
	if err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, account)
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	accountNumber := chi.URLParam(r, "accountNumber")
	account, err := h.accounts.Get(r.Context(), accountNumber)
	if err != nil {
		ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, account)
}

func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountNumber := chi.URLParam(r, "accountNumber")
	var req model.UpdateAccountRequest
	if err := DecodeJSON(r, &req); err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	account, err := h.accounts.Update(r.Context(), accountNumber, req)
	if err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, account)
}

func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountNumber := chi.URLParam(r, "accountNumber")
	if err := h.accounts.Delete(r.Context(), accountNumber); err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"message": "account deleted"})
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	params := model.ListParams{
		Limit:  QueryInt(r, "limit", 20),
		Offset: QueryInt(r, "offset", 0),
	}

	accounts, total, err := h.accounts.List(r.Context(), params)
	if err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{
		"accounts": accounts,
		"total":    total,
		"limit":    params.Limit,
		"offset":   params.Offset,
	})
}

func (h *AccountHandler) ListByCustomer(w http.ResponseWriter, r *http.Request) {
	customerNumber := chi.URLParam(r, "customerNumber")
	accounts, err := h.accounts.ListByCustomer(r.Context(), customerNumber)
	if err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{
		"accounts": accounts,
		"total":    len(accounts),
	})
}

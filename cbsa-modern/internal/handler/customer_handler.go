package handler

import (
	"net/http"

	"github.com/cicsdev/cbsa-modern/internal/model"
	"github.com/cicsdev/cbsa-modern/internal/service"
	"github.com/go-chi/chi/v5"
)

type CustomerHandler struct {
	customers *service.CustomerService
	credit    *service.CreditService
}

func NewCustomerHandler(customers *service.CustomerService, credit *service.CreditService) *CustomerHandler {
	return &CustomerHandler{customers: customers, credit: credit}
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCustomerRequest
	if err := DecodeJSON(r, &req); err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.FirstName == "" || req.LastName == "" || req.DateOfBirth == "" {
		ErrorJSON(w, http.StatusBadRequest, "firstName, lastName, and dateOfBirth are required")
		return
	}

	creditScore, err := h.credit.CheckCredit(r.Context(), "")
	if err != nil {
		ErrorJSON(w, http.StatusInternalServerError, "credit check failed: "+err.Error())
		return
	}

	customer, err := h.customers.Create(r.Context(), req, creditScore)
	if err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	customerNumber := chi.URLParam(r, "customerNumber")
	customer, err := h.customers.Get(r.Context(), customerNumber)
	if err != nil {
		ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	customerNumber := chi.URLParam(r, "customerNumber")
	var req model.UpdateCustomerRequest
	if err := DecodeJSON(r, &req); err != nil {
		ErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	customer, err := h.customers.Update(r.Context(), customerNumber, req)
	if err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	customerNumber := chi.URLParam(r, "customerNumber")
	if err := h.customers.Delete(r.Context(), customerNumber); err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"message": "customer deleted"})
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	params := model.ListParams{
		Limit:  QueryInt(r, "limit", 20),
		Offset: QueryInt(r, "offset", 0),
	}

	customers, total, err := h.customers.List(r.Context(), params)
	if err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]interface{}{
		"customers": customers,
		"total":     total,
		"limit":     params.Limit,
		"offset":    params.Offset,
	})
}

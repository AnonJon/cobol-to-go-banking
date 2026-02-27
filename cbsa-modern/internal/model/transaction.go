package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProcessedTransaction struct {
	ID            int64           `json:"id" db:"id"`
	SortCode      string          `json:"sortCode" db:"sort_code"`
	AccountNumber string          `json:"accountNumber" db:"account_number"`
	TxnDate       time.Time       `json:"txnDate" db:"txn_date"`
	TxnTime       time.Time       `json:"txnTime" db:"txn_time"`
	Reference     string          `json:"reference" db:"reference"`
	TxnType       string          `json:"txnType" db:"txn_type"`
	Description   string          `json:"description" db:"description"`
	Amount        decimal.Decimal `json:"amount" db:"amount"`
	CreatedAt     time.Time       `json:"createdAt" db:"created_at"`
}

type DebitCreditRequest struct {
	AccountNumber string `json:"accountNumber"`
	Amount        string `json:"amount"`
	Description   string `json:"description"`
}

type TransferRequest struct {
	FromAccountNumber string `json:"fromAccountNumber"`
	ToAccountNumber   string `json:"toAccountNumber"`
	Amount            string `json:"amount"`
	Description       string `json:"description"`
}

type ListParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

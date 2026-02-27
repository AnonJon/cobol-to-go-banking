package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID               int64           `json:"id" db:"id"`
	AccountNumber    string          `json:"accountNumber" db:"account_number"`
	SortCode         string          `json:"sortCode" db:"sort_code"`
	CustomerNumber   string          `json:"customerNumber" db:"customer_number"`
	AccountType      string          `json:"accountType" db:"account_type"`
	InterestRate     decimal.Decimal `json:"interestRate" db:"interest_rate"`
	OpenedDate       time.Time       `json:"openedDate" db:"opened_date"`
	OverdraftLimit   int             `json:"overdraftLimit" db:"overdraft_limit"`
	LastStatement    time.Time       `json:"lastStatement" db:"last_statement"`
	NextStatement    time.Time       `json:"nextStatement" db:"next_statement"`
	AvailableBalance decimal.Decimal `json:"availableBalance" db:"available_balance"`
	ActualBalance    decimal.Decimal `json:"actualBalance" db:"actual_balance"`
	CreatedAt        time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time       `json:"updatedAt" db:"updated_at"`
}

type CreateAccountRequest struct {
	CustomerNumber string `json:"customerNumber"`
	AccountType    string `json:"accountType"`
	InterestRate   string `json:"interestRate"`
	OverdraftLimit int    `json:"overdraftLimit"`
}

type UpdateAccountRequest struct {
	AccountType    *string `json:"accountType,omitempty"`
	InterestRate   *string `json:"interestRate,omitempty"`
	OverdraftLimit *int    `json:"overdraftLimit,omitempty"`
}

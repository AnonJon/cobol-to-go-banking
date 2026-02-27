package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// BenchmarkAccountMarshal measures JSON serialization speed for the Account
// model. This is the hot path for every account inquiry API call, replacing
// the COBOL INQACC.cbl -> z/OS Connect JSON transformation pipeline.
func BenchmarkAccountMarshal(b *testing.B) {
	acct := Account{
		ID:               1,
		AccountNumber:    "00000001",
		SortCode:         "987654",
		CustomerNumber:   "0000000001",
		AccountType:      "CURRENT",
		InterestRate:     decimal.NewFromFloat(2.25),
		OpenedDate:       time.Now(),
		OverdraftLimit:   500,
		LastStatement:    time.Now(),
		NextStatement:    time.Now().AddDate(0, 1, 0),
		AvailableBalance: decimal.NewFromFloat(12345.67),
		ActualBalance:    decimal.NewFromFloat(12345.67),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(acct)
	}
}

// BenchmarkDecimalArithmetic measures the performance of financial calculations
// using shopspring/decimal, replacing COBOL PIC S9(10)V99 COMP-3 arithmetic.
func BenchmarkDecimalArithmetic(b *testing.B) {
	balance := decimal.NewFromFloat(50000.00)
	amount := decimal.NewFromFloat(123.45)
	overdraft := decimal.NewFromInt(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newBal := balance.Sub(amount)
		_ = newBal.LessThan(overdraft.Neg())
		balance = newBal.Add(amount)
	}
}

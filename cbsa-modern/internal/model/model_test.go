package model

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
)

func TestAccount_JSONRoundTrip(t *testing.T) {
	acct := Account{
		ID:               1,
		AccountNumber:    "00000001",
		SortCode:         "987654",
		CustomerNumber:   "0000000001",
		AccountType:      "CURRENT",
		InterestRate:     decimal.NewFromFloat(2.25),
		OverdraftLimit:   500,
		AvailableBalance: decimal.NewFromFloat(1234.56),
		ActualBalance:    decimal.NewFromFloat(1234.56),
	}

	data, err := json.Marshal(acct)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Account
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !decoded.AvailableBalance.Equal(acct.AvailableBalance) {
		t.Errorf("balance mismatch: got %s, want %s",
			decoded.AvailableBalance.StringFixed(2),
			acct.AvailableBalance.StringFixed(2))
	}

	if decoded.AccountNumber != acct.AccountNumber {
		t.Errorf("account number mismatch: got %s, want %s",
			decoded.AccountNumber, acct.AccountNumber)
	}
}

func TestDecimalPrecision(t *testing.T) {
	// Verify that shopspring/decimal preserves COBOL PIC S9(10)V99 precision
	tests := []struct {
		name   string
		a, b   string
		op     string
		expect string
	}{
		{"add", "1234567890.12", "0.01", "add", "1234567890.13"},
		{"subtract", "100.00", "50.50", "sub", "49.50"},
		{"large balance", "9999999999.99", "0.01", "add", "10000000000.00"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, _ := decimal.NewFromString(tc.a)
			b, _ := decimal.NewFromString(tc.b)
			expected, _ := decimal.NewFromString(tc.expect)

			var result decimal.Decimal
			switch tc.op {
			case "add":
				result = a.Add(b)
			case "sub":
				result = a.Sub(b)
			}

			if !result.Equal(expected) {
				t.Errorf("%s %s %s = %s, want %s",
					tc.a, tc.op, tc.b, result.StringFixed(2), expected.StringFixed(2))
			}
		})
	}
}

func TestListParams_Defaults(t *testing.T) {
	params := ListParams{}
	if params.Limit != 0 {
		t.Errorf("expected zero value for Limit, got %d", params.Limit)
	}
	if params.Offset != 0 {
		t.Errorf("expected zero value for Offset, got %d", params.Offset)
	}
}

func TestProcessedTransaction_JSON(t *testing.T) {
	txn := ProcessedTransaction{
		TxnType:     "DEB",
		Amount:      decimal.NewFromFloat(50.00),
		Description: "Cash withdrawal",
	}

	data, err := json.Marshal(txn)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded ProcessedTransaction
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.TxnType != "DEB" {
		t.Errorf("txn type mismatch: got %s, want DEB", decoded.TxnType)
	}
	if !decoded.Amount.Equal(txn.Amount) {
		t.Errorf("amount mismatch: got %s, want %s",
			decoded.Amount.StringFixed(2), txn.Amount.StringFixed(2))
	}
}

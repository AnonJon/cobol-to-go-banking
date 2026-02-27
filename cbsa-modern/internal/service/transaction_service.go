package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cicsdev/cbsa-modern/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

const maxRetries = 3

type TransactionService struct {
	db       *sqlx.DB
	sortCode string
}

func NewTransactionService(db *sqlx.DB, sortCode string) *TransactionService {
	return &TransactionService{db: db, sortCode: sortCode}
}

// Debit replaces the debit path of DBCRFUN.cbl — retrieves the account with
// SELECT FOR UPDATE (replacing CICS ENQ), validates sufficient funds against
// the overdraft limit, updates both balances, and logs to processed_transactions.
func (s *TransactionService) Debit(ctx context.Context, req model.DebitCreditRequest) (*model.Account, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("invalid amount")
	}

	var result model.Account
	err = s.withRetry(ctx, func(tx *sqlx.Tx) error {
		var acct model.Account
		err := tx.QueryRowxContext(ctx,
			`SELECT * FROM accounts WHERE account_number=$1 AND sort_code=$2 FOR UPDATE`,
			req.AccountNumber, s.sortCode).StructScan(&acct)
		if err != nil {
			return fmt.Errorf("account %s not found: %w", req.AccountNumber, err)
		}

		overdraftDec := decimal.NewFromInt(int64(acct.OverdraftLimit))
		if acct.AvailableBalance.Sub(amount).LessThan(overdraftDec.Neg()) {
			return fmt.Errorf("insufficient funds: available %s, overdraft limit %s",
				acct.AvailableBalance.StringFixed(2), overdraftDec.StringFixed(2))
		}

		acct.AvailableBalance = acct.AvailableBalance.Sub(amount)
		acct.ActualBalance = acct.ActualBalance.Sub(amount)

		err = tx.QueryRowContext(ctx,
			`UPDATE accounts SET available_balance=$1, actual_balance=$2, updated_at=now()
			 WHERE account_number=$3 AND sort_code=$4 RETURNING updated_at`,
			acct.AvailableBalance, acct.ActualBalance, req.AccountNumber, s.sortCode,
		).Scan(&acct.UpdatedAt)
		if err != nil {
			return fmt.Errorf("updating balance: %w", err)
		}

		_, err = tx.ExecContext(ctx,
			`INSERT INTO processed_transactions (sort_code, account_number, txn_date, txn_time, reference, txn_type, description, amount)
			 VALUES ($1, $2, CURRENT_DATE, CURRENT_TIME, $3, 'DEB', $4, $5)`,
			s.sortCode, req.AccountNumber, req.AccountNumber, req.Description, amount)
		if err != nil {
			return fmt.Errorf("logging transaction: %w", err)
		}

		result = acct
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Credit replaces the credit path of DBCRFUN.cbl — adds funds to an account
// and logs the transaction.
func (s *TransactionService) Credit(ctx context.Context, req model.DebitCreditRequest) (*model.Account, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("invalid amount")
	}

	var result model.Account
	err = s.withRetry(ctx, func(tx *sqlx.Tx) error {
		var acct model.Account
		err := tx.QueryRowxContext(ctx,
			`SELECT * FROM accounts WHERE account_number=$1 AND sort_code=$2 FOR UPDATE`,
			req.AccountNumber, s.sortCode).StructScan(&acct)
		if err != nil {
			return fmt.Errorf("account %s not found: %w", req.AccountNumber, err)
		}

		acct.AvailableBalance = acct.AvailableBalance.Add(amount)
		acct.ActualBalance = acct.ActualBalance.Add(amount)

		err = tx.QueryRowContext(ctx,
			`UPDATE accounts SET available_balance=$1, actual_balance=$2, updated_at=now()
			 WHERE account_number=$3 AND sort_code=$4 RETURNING updated_at`,
			acct.AvailableBalance, acct.ActualBalance, req.AccountNumber, s.sortCode,
		).Scan(&acct.UpdatedAt)
		if err != nil {
			return fmt.Errorf("updating balance: %w", err)
		}

		_, err = tx.ExecContext(ctx,
			`INSERT INTO processed_transactions (sort_code, account_number, txn_date, txn_time, reference, txn_type, description, amount)
			 VALUES ($1, $2, CURRENT_DATE, CURRENT_TIME, $3, 'CRE', $4, $5)`,
			s.sortCode, req.AccountNumber, req.AccountNumber, req.Description, amount)
		if err != nil {
			return fmt.Errorf("logging transaction: %w", err)
		}

		result = acct
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Transfer replaces XFRFUN.cbl — moves funds between two accounts in a single
// transaction. Both accounts are locked with SELECT FOR UPDATE (replacing
// CICS ENQ/DEQ) and the transfer is logged for both sides. Includes retry
// logic replacing COBOL's SQLCODE -911/-913 deadlock handling.
func (s *TransactionService) Transfer(ctx context.Context, req model.TransferRequest) error {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("invalid amount")
	}
	if req.FromAccountNumber == req.ToAccountNumber {
		return fmt.Errorf("cannot transfer to the same account")
	}

	return s.withRetry(ctx, func(tx *sqlx.Tx) error {
		// Lock both accounts in a consistent order to avoid deadlocks
		first, second := req.FromAccountNumber, req.ToAccountNumber
		if first > second {
			first, second = second, first
		}

		var acct1, acct2 model.Account
		err := tx.QueryRowxContext(ctx,
			`SELECT * FROM accounts WHERE account_number=$1 AND sort_code=$2 FOR UPDATE`,
			first, s.sortCode).StructScan(&acct1)
		if err != nil {
			return fmt.Errorf("account %s not found: %w", first, err)
		}
		err = tx.QueryRowxContext(ctx,
			`SELECT * FROM accounts WHERE account_number=$1 AND sort_code=$2 FOR UPDATE`,
			second, s.sortCode).StructScan(&acct2)
		if err != nil {
			return fmt.Errorf("account %s not found: %w", second, err)
		}

		var fromAcct, toAcct *model.Account
		if acct1.AccountNumber == req.FromAccountNumber {
			fromAcct, toAcct = &acct1, &acct2
		} else {
			fromAcct, toAcct = &acct2, &acct1
		}

		overdraftDec := decimal.NewFromInt(int64(fromAcct.OverdraftLimit))
		if fromAcct.AvailableBalance.Sub(amount).LessThan(overdraftDec.Neg()) {
			return fmt.Errorf("insufficient funds in account %s", req.FromAccountNumber)
		}

		fromAcct.AvailableBalance = fromAcct.AvailableBalance.Sub(amount)
		fromAcct.ActualBalance = fromAcct.ActualBalance.Sub(amount)
		toAcct.AvailableBalance = toAcct.AvailableBalance.Add(amount)
		toAcct.ActualBalance = toAcct.ActualBalance.Add(amount)

		_, err = tx.ExecContext(ctx,
			`UPDATE accounts SET available_balance=$1, actual_balance=$2, updated_at=now()
			 WHERE account_number=$3 AND sort_code=$4`,
			fromAcct.AvailableBalance, fromAcct.ActualBalance, req.FromAccountNumber, s.sortCode)
		if err != nil {
			return fmt.Errorf("updating from account: %w", err)
		}

		_, err = tx.ExecContext(ctx,
			`UPDATE accounts SET available_balance=$1, actual_balance=$2, updated_at=now()
			 WHERE account_number=$3 AND sort_code=$4`,
			toAcct.AvailableBalance, toAcct.ActualBalance, req.ToAccountNumber, s.sortCode)
		if err != nil {
			return fmt.Errorf("updating to account: %w", err)
		}

		desc := req.Description
		if desc == "" {
			desc = fmt.Sprintf("Transfer %s -> %s", req.FromAccountNumber, req.ToAccountNumber)
		}

		for _, entry := range []struct {
			acctNum string
			txnType string
		}{
			{req.FromAccountNumber, "TFR"},
			{req.ToAccountNumber, "TFR"},
		} {
			_, err = tx.ExecContext(ctx,
				`INSERT INTO processed_transactions (sort_code, account_number, txn_date, txn_time, reference, txn_type, description, amount)
				 VALUES ($1, $2, CURRENT_DATE, CURRENT_TIME, $3, $4, $5, $6)`,
				s.sortCode, entry.acctNum, entry.acctNum, entry.txnType, desc, amount)
			if err != nil {
				return fmt.Errorf("logging transfer: %w", err)
			}
		}

		return nil
	})
}

// List retrieves processed transactions with pagination.
func (s *TransactionService) List(ctx context.Context, params model.ListParams) ([]model.ProcessedTransaction, int, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}

	var total int
	err := s.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM processed_transactions WHERE sort_code=$1`, s.sortCode)
	if err != nil {
		return nil, 0, fmt.Errorf("counting transactions: %w", err)
	}

	var txns []model.ProcessedTransaction
	err = s.db.SelectContext(ctx, &txns,
		`SELECT * FROM processed_transactions WHERE sort_code=$1 ORDER BY id DESC LIMIT $2 OFFSET $3`,
		s.sortCode, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("listing transactions: %w", err)
	}
	return txns, total, nil
}

// withRetry wraps transactional work with retry logic, replacing the COBOL
// SQLCODE -911/-913 deadlock/timeout handling in XFRFUN.cbl and DBCRFUN.cbl.
func (s *TransactionService) withRetry(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	for attempt := 0; attempt < maxRetries; attempt++ {
		tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			return fmt.Errorf("begin transaction: %w", err)
		}

		if err := fn(tx); err != nil {
			tx.Rollback()
			if attempt < maxRetries-1 && isRetryable(err) {
				continue
			}
			return err
		}

		if err := tx.Commit(); err != nil {
			if attempt < maxRetries-1 && isRetryable(err) {
				continue
			}
			return fmt.Errorf("committing: %w", err)
		}
		return nil
	}
	return fmt.Errorf("transaction failed after %d retries", maxRetries)
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, keyword := range []string{"deadlock", "serialization", "could not serialize"} {
		if contains(msg, keyword) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

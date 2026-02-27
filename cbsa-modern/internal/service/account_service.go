package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cicsdev/cbsa-modern/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type AccountService struct {
	db       *sqlx.DB
	sortCode string
}

func NewAccountService(db *sqlx.DB, sortCode string) *AccountService {
	return &AccountService{db: db, sortCode: sortCode}
}

// Create replaces CREACC.cbl — validates the customer exists, obtains a new
// account number from a sequence (replacing CICS Named Counter + ENQ/DEQ),
// inserts the account, and logs to processed_transactions.
func (s *AccountService) Create(ctx context.Context, req model.CreateAccountRequest) (*model.Account, error) {
	interestRate, err := decimal.NewFromString(req.InterestRate)
	if err != nil {
		return nil, fmt.Errorf("invalid interest rate: %w", err)
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM customers WHERE customer_number=$1 AND sort_code=$2)`,
		req.CustomerNumber, s.sortCode).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("checking customer: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("customer %s not found", req.CustomerNumber)
	}

	var seq int64
	if err := tx.QueryRowContext(ctx, "SELECT nextval('account_number_seq')").Scan(&seq); err != nil {
		return nil, fmt.Errorf("generating account number: %w", err)
	}
	accountNumber := fmt.Sprintf("%08d", seq)

	now := time.Now()
	nextStatement := now.AddDate(0, 1, 0)

	acct := model.Account{
		AccountNumber:    accountNumber,
		SortCode:         s.sortCode,
		CustomerNumber:   req.CustomerNumber,
		AccountType:      req.AccountType,
		InterestRate:     interestRate,
		OpenedDate:       now,
		OverdraftLimit:   req.OverdraftLimit,
		LastStatement:    now,
		NextStatement:    nextStatement,
		AvailableBalance: decimal.Zero,
		ActualBalance:    decimal.Zero,
	}

	query := `INSERT INTO accounts (account_number, sort_code, customer_number, account_type,
		interest_rate, opened_date, overdraft_limit, last_statement, next_statement,
		available_balance, actual_balance)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query,
		acct.AccountNumber, acct.SortCode, acct.CustomerNumber, acct.AccountType,
		acct.InterestRate, acct.OpenedDate, acct.OverdraftLimit,
		acct.LastStatement, acct.NextStatement,
		acct.AvailableBalance, acct.ActualBalance,
	).Scan(&acct.ID, &acct.CreatedAt, &acct.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("inserting account: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO processed_transactions (sort_code, account_number, txn_date, txn_time, reference, txn_type, description, amount)
		 VALUES ($1, $2, CURRENT_DATE, CURRENT_TIME, $3, 'OAC', $4, 0)`,
		s.sortCode, accountNumber, accountNumber, "Account "+accountNumber+" opened")
	if err != nil {
		return nil, fmt.Errorf("logging transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing: %w", err)
	}
	return &acct, nil
}

// Get replaces INQACC.cbl — retrieves a single account by number
// (replacing DB2 SELECT with cursor).
func (s *AccountService) Get(ctx context.Context, accountNumber string) (*model.Account, error) {
	var a model.Account
	err := s.db.GetContext(ctx, &a,
		`SELECT * FROM accounts WHERE account_number=$1 AND sort_code=$2`,
		accountNumber, s.sortCode)
	if err != nil {
		return nil, fmt.Errorf("account %s not found: %w", accountNumber, err)
	}
	return &a, nil
}

// ListByCustomer replaces INQACCCU.cbl — retrieves all accounts for a customer
// (replacing DB2 cursor iteration + INQCUST validation).
func (s *AccountService) ListByCustomer(ctx context.Context, customerNumber string) ([]model.Account, error) {
	var accounts []model.Account
	err := s.db.SelectContext(ctx, &accounts,
		`SELECT * FROM accounts WHERE customer_number=$1 AND sort_code=$2 ORDER BY account_number`,
		customerNumber, s.sortCode)
	if err != nil {
		return nil, fmt.Errorf("listing accounts for customer %s: %w", customerNumber, err)
	}
	return accounts, nil
}

// Update replaces UPDACC.cbl — updates non-financial account fields
// (replacing DB2 SELECT + UPDATE with COBOL field validation).
func (s *AccountService) Update(ctx context.Context, accountNumber string, req model.UpdateAccountRequest) (*model.Account, error) {
	existing, err := s.Get(ctx, accountNumber)
	if err != nil {
		return nil, err
	}

	if req.AccountType != nil {
		existing.AccountType = *req.AccountType
	}
	if req.InterestRate != nil {
		rate, err := decimal.NewFromString(*req.InterestRate)
		if err != nil {
			return nil, fmt.Errorf("invalid interest rate: %w", err)
		}
		existing.InterestRate = rate
	}
	if req.OverdraftLimit != nil {
		existing.OverdraftLimit = *req.OverdraftLimit
	}

	query := `UPDATE accounts SET account_type=$1, interest_rate=$2, overdraft_limit=$3, updated_at=now()
		WHERE account_number=$4 AND sort_code=$5
		RETURNING updated_at`

	err = s.db.QueryRowContext(ctx, query,
		existing.AccountType, existing.InterestRate, existing.OverdraftLimit,
		accountNumber, s.sortCode,
	).Scan(&existing.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updating account: %w", err)
	}
	return existing, nil
}

// Delete replaces DELACC.cbl — deletes an account and logs to processed_transactions
// (replacing DB2 DELETE + PROCTRAN INSERT).
func (s *AccountService) Delete(ctx context.Context, accountNumber string) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO processed_transactions (sort_code, account_number, txn_date, txn_time, reference, txn_type, description, amount)
		 VALUES ($1, $2, CURRENT_DATE, CURRENT_TIME, $3, 'DAC', $4, 0)`,
		s.sortCode, accountNumber, accountNumber, "Account "+accountNumber+" deleted")
	if err != nil {
		return fmt.Errorf("logging transaction: %w", err)
	}

	result, err := tx.ExecContext(ctx,
		`DELETE FROM accounts WHERE account_number=$1 AND sort_code=$2`,
		accountNumber, s.sortCode)
	if err != nil {
		return fmt.Errorf("deleting account: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("account %s not found", accountNumber)
	}

	return tx.Commit()
}

// List supports paginated account listing.
func (s *AccountService) List(ctx context.Context, params model.ListParams) ([]model.Account, int, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}

	var total int
	err := s.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM accounts WHERE sort_code=$1`, s.sortCode)
	if err != nil {
		return nil, 0, fmt.Errorf("counting accounts: %w", err)
	}

	var accounts []model.Account
	err = s.db.SelectContext(ctx, &accounts,
		`SELECT * FROM accounts WHERE sort_code=$1 ORDER BY account_number LIMIT $2 OFFSET $3`,
		s.sortCode, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("listing accounts: %w", err)
	}
	return accounts, total, nil
}

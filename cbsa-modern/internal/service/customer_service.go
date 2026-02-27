package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cicsdev/cbsa-modern/internal/model"
	"github.com/jmoiron/sqlx"
)

type CustomerService struct {
	db       *sqlx.DB
	sortCode string
}

func NewCustomerService(db *sqlx.DB, sortCode string) *CustomerService {
	return &CustomerService{db: db, sortCode: sortCode}
}

// Create replaces CRECUST.cbl — creates a customer record with a sequenced
// customer number (replacing CICS Named Counter) and logs the creation in
// processed_transactions (replacing PROCTRAN INSERT).
func (s *CustomerService) Create(ctx context.Context, req model.CreateCustomerRequest, creditScore int) (*model.Customer, error) {
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("invalid date of birth: %w", err)
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var seq int64
	if err := tx.QueryRowContext(ctx, "SELECT nextval('customer_number_seq')").Scan(&seq); err != nil {
		return nil, fmt.Errorf("generating customer number: %w", err)
	}
	customerNumber := fmt.Sprintf("%010d", seq)

	now := time.Now()
	reviewDate := now.AddDate(0, 6, 0)

	customer := model.Customer{
		CustomerNumber: customerNumber,
		SortCode:       s.sortCode,
		Title:          req.Title,
		FirstName:      req.FirstName,
		LastName:        req.LastName,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   req.AddressLine2,
		AddressLine3:   req.AddressLine3,
		DateOfBirth:    dob,
		CreditScore:    creditScore,
		ReviewDate:     &reviewDate,
	}

	query := `INSERT INTO customers (customer_number, sort_code, title, first_name, last_name,
		address_line1, address_line2, address_line3, date_of_birth, credit_score, review_date)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query,
		customer.CustomerNumber, customer.SortCode, customer.Title,
		customer.FirstName, customer.LastName,
		customer.AddressLine1, customer.AddressLine2, customer.AddressLine3,
		customer.DateOfBirth, customer.CreditScore, customer.ReviewDate,
	).Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("inserting customer: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO processed_transactions (sort_code, account_number, txn_date, txn_time, reference, txn_type, description, amount)
		 VALUES ($1, '00000000', CURRENT_DATE, CURRENT_TIME, $2, 'CCS', $3, 0)`,
		s.sortCode, customerNumber, "Customer "+customerNumber+" created")
	if err != nil {
		return nil, fmt.Errorf("logging transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing: %w", err)
	}
	return &customer, nil
}

// Get replaces INQCUST.cbl — reads a customer by number (replacing VSAM READ FILE).
func (s *CustomerService) Get(ctx context.Context, customerNumber string) (*model.Customer, error) {
	var c model.Customer
	err := s.db.GetContext(ctx, &c,
		`SELECT * FROM customers WHERE customer_number=$1 AND sort_code=$2`,
		customerNumber, s.sortCode)
	if err != nil {
		return nil, fmt.Errorf("customer %s not found: %w", customerNumber, err)
	}
	return &c, nil
}

// Update replaces UPDCUST.cbl — updates a customer record (replacing VSAM
// READ FILE + REWRITE FILE with name/address validation).
func (s *CustomerService) Update(ctx context.Context, customerNumber string, req model.UpdateCustomerRequest) (*model.Customer, error) {
	existing, err := s.Get(ctx, customerNumber)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.FirstName != nil {
		existing.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		existing.LastName = *req.LastName
	}
	if req.AddressLine1 != nil {
		existing.AddressLine1 = *req.AddressLine1
	}
	if req.AddressLine2 != nil {
		existing.AddressLine2 = *req.AddressLine2
	}
	if req.AddressLine3 != nil {
		existing.AddressLine3 = *req.AddressLine3
	}

	query := `UPDATE customers SET title=$1, first_name=$2, last_name=$3,
		address_line1=$4, address_line2=$5, address_line3=$6, updated_at=now()
		WHERE customer_number=$7 AND sort_code=$8
		RETURNING updated_at`

	err = s.db.QueryRowContext(ctx, query,
		existing.Title, existing.FirstName, existing.LastName,
		existing.AddressLine1, existing.AddressLine2, existing.AddressLine3,
		customerNumber, s.sortCode,
	).Scan(&existing.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updating customer: %w", err)
	}
	return existing, nil
}

// Delete replaces DELCUS.cbl — deletes a customer and cascades to all their
// accounts (replacing the COBOL loop through INQACCCU + DELACC + VSAM DELETE).
func (s *CustomerService) Delete(ctx context.Context, customerNumber string) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO processed_transactions (sort_code, account_number, txn_date, txn_time, reference, txn_type, description, amount)
		 VALUES ($1, '00000000', CURRENT_DATE, CURRENT_TIME, $2, 'DCS', $3, 0)`,
		s.sortCode, customerNumber, "Customer "+customerNumber+" deleted")
	if err != nil {
		return fmt.Errorf("logging transaction: %w", err)
	}

	// ON DELETE CASCADE handles account cleanup
	result, err := tx.ExecContext(ctx,
		`DELETE FROM customers WHERE customer_number=$1 AND sort_code=$2`,
		customerNumber, s.sortCode)
	if err != nil {
		return fmt.Errorf("deleting customer: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("customer %s not found", customerNumber)
	}

	return tx.Commit()
}

// List replaces the Java CustomerResource list logic, supporting pagination.
func (s *CustomerService) List(ctx context.Context, params model.ListParams) ([]model.Customer, int, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}

	var total int
	err := s.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM customers WHERE sort_code=$1`, s.sortCode)
	if err != nil {
		return nil, 0, fmt.Errorf("counting customers: %w", err)
	}

	var customers []model.Customer
	err = s.db.SelectContext(ctx, &customers,
		`SELECT * FROM customers WHERE sort_code=$1 ORDER BY customer_number LIMIT $2 OFFSET $3`,
		s.sortCode, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("listing customers: %w", err)
	}
	return customers, total, nil
}

package model

import "time"

type Customer struct {
	ID             int64     `json:"id" db:"id"`
	CustomerNumber string    `json:"customerNumber" db:"customer_number"`
	SortCode       string    `json:"sortCode" db:"sort_code"`
	Title          string    `json:"title" db:"title"`
	FirstName      string    `json:"firstName" db:"first_name"`
	LastName       string    `json:"lastName" db:"last_name"`
	AddressLine1   string    `json:"addressLine1" db:"address_line1"`
	AddressLine2   string    `json:"addressLine2" db:"address_line2"`
	AddressLine3   string    `json:"addressLine3" db:"address_line3"`
	DateOfBirth    time.Time `json:"dateOfBirth" db:"date_of_birth"`
	CreditScore    int       `json:"creditScore" db:"credit_score"`
	ReviewDate     *time.Time `json:"reviewDate,omitempty" db:"review_date"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateCustomerRequest struct {
	Title        string `json:"title"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	DateOfBirth  string `json:"dateOfBirth"` // YYYY-MM-DD
}

type UpdateCustomerRequest struct {
	Title        *string `json:"title,omitempty"`
	FirstName    *string `json:"firstName,omitempty"`
	LastName     *string `json:"lastName,omitempty"`
	AddressLine1 *string `json:"addressLine1,omitempty"`
	AddressLine2 *string `json:"addressLine2,omitempty"`
	AddressLine3 *string `json:"addressLine3,omitempty"`
}

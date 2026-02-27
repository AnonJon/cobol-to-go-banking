package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/cicsdev/cbsa-modern/internal/config"
	"github.com/cicsdev/cbsa-modern/internal/database"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// seed replaces BANKDATA.cbl — the batch COBOL program that populated the
// VSAM CUSTOMER file and Db2 ACCOUNT table with randomized test data.
// The COBOL version used JCL PARM for start/end keys and CEEGMT/CEEDATM
// for timestamps. This Go version generates the same style of data.

var (
	titles    = []string{"Mr", "Mrs", "Miss", "Ms", "Dr", "Prof"}
	firstNames = []string{
		"James", "Mary", "Robert", "Patricia", "John", "Jennifer", "Michael", "Linda",
		"David", "Elizabeth", "William", "Barbara", "Richard", "Susan", "Joseph", "Jessica",
		"Thomas", "Sarah", "Charles", "Karen", "Christopher", "Lisa", "Daniel", "Nancy",
		"Matthew", "Betty", "Anthony", "Margaret", "Mark", "Sandra",
	}
	lastNames = []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
		"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson",
		"Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson",
		"White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson",
	}
	streets = []string{
		"Oak Avenue", "Elm Street", "Maple Drive", "Pine Road", "Cedar Lane",
		"Birch Way", "Willow Court", "Ash Boulevard", "Cherry Place", "Walnut Street",
		"High Street", "Station Road", "Church Lane", "Mill Road", "Park Avenue",
	}
	towns = []string{
		"London", "Birmingham", "Manchester", "Leeds", "Liverpool",
		"Bristol", "Sheffield", "Edinburgh", "Glasgow", "Cardiff",
		"Oxford", "Cambridge", "Brighton", "Bath", "York",
	}
	accountTypes = []string{"ISA", "SAVING", "CURRENT", "LOAN", "MORTGAGE"}
	interestRates = []string{"1.50", "2.25", "3.00", "0.50", "4.75", "1.00", "2.50"}
	overdraftLimits = []int{0, 100, 250, 500, 1000, 2500, 5000}
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = "/migrations"
	}
	if err := database.RunMigrations(db, migrationsDir); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	numCustomers := 100
	maxAccountsPerCustomer := 5

	slog.Info("seeding data", "customers", numCustomers, "maxAccountsPer", maxAccountsPerCustomer)

	for i := 1; i <= numCustomers; i++ {
		custNum := fmt.Sprintf("%010d", i)
		if err := seedCustomer(db, cfg.SortCode, custNum); err != nil {
			slog.Error("failed to seed customer", "number", custNum, "error", err)
			continue
		}

		numAccounts := rand.Intn(maxAccountsPerCustomer) + 1
		for j := 0; j < numAccounts; j++ {
			acctNum := fmt.Sprintf("%08d", (i-1)*maxAccountsPerCustomer+j+1)
			if err := seedAccount(db, cfg.SortCode, custNum, acctNum); err != nil {
				slog.Error("failed to seed account", "number", acctNum, "error", err)
			}
		}
	}

	slog.Info("seeding complete")
}

func seedCustomer(db *sqlx.DB, sortCode, customerNumber string) error {
	dob := time.Date(
		1950+rand.Intn(50),
		time.Month(rand.Intn(12)+1),
		rand.Intn(28)+1,
		0, 0, 0, 0, time.UTC,
	)
	reviewDate := time.Now().AddDate(0, 6, 0)

	_, err := db.Exec(
		`INSERT INTO customers (customer_number, sort_code, title, first_name, last_name,
		 address_line1, address_line2, address_line3, date_of_birth, credit_score, review_date)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		 ON CONFLICT (customer_number) DO NOTHING`,
		customerNumber,
		sortCode,
		titles[rand.Intn(len(titles))],
		firstNames[rand.Intn(len(firstNames))],
		lastNames[rand.Intn(len(lastNames))],
		fmt.Sprintf("%d %s", rand.Intn(200)+1, streets[rand.Intn(len(streets))]),
		towns[rand.Intn(len(towns))],
		fmt.Sprintf("UK, %s%d %s%s",
			string(rune('A'+rand.Intn(26))),
			rand.Intn(9)+1,
			string(rune('0'+rand.Intn(10))),
			string(rune('A'+rand.Intn(26)))),
		dob,
		rand.Intn(999)+1,
		reviewDate,
	)
	return err
}

func seedAccount(db *sqlx.DB, sortCode, customerNumber, accountNumber string) error {
	openDate := time.Now().AddDate(0, -rand.Intn(60), -rand.Intn(28))
	lastStmt := openDate.AddDate(0, rand.Intn(3), 0)
	nextStmt := lastStmt.AddDate(0, 1, 0)

	rate, _ := decimal.NewFromString(interestRates[rand.Intn(len(interestRates))])
	balance := decimal.NewFromFloat(float64(rand.Intn(100000)) / 100.0)

	_, err := db.Exec(
		`INSERT INTO accounts (account_number, sort_code, customer_number, account_type,
		 interest_rate, opened_date, overdraft_limit, last_statement, next_statement,
		 available_balance, actual_balance)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		 ON CONFLICT (account_number) DO NOTHING`,
		accountNumber,
		sortCode,
		customerNumber,
		accountTypes[rand.Intn(len(accountTypes))],
		rate,
		openDate,
		overdraftLimits[rand.Intn(len(overdraftLimits))],
		lastStmt,
		nextStmt,
		balance,
		balance,
	)
	return err
}

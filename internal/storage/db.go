package storage

import (
	"banking-ledger-service/internal/models"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// InitDB initializes the PostgreSQL database connection
func InitDB() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost" // Default for local development
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "user" // Default user
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password" // Default password
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "banking_ledger" // Default database name
	}

	connStr := "postgres://" + user + ":" + password + "@" + host + ":5432/" + dbName + "?sslmode=disable"

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	DB = dbpool
	log.Println("Connected to PostgreSQL")
}

// CreateAccount inserts a new account while ensuring uniqueness
func CreateAccount(name string, balance float64) (int, error) {
	var id int
	tx, err := DB.Begin(context.Background())
	if err != nil {
		return 0, err
	}

	// Rollback transaction if any error occurs
	defer tx.Rollback(context.Background())

	// Insert new account
	err = tx.QueryRow(context.Background(), "INSERT INTO accounts (name, balance) VALUES ($1, $2) RETURNING id", name, balance).Scan(&id)
	if err != nil {
		return 0, err
	}

	// Commit transaction
	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	// Add transaction record
	err = AddTransaction(id, balance, "account_creation")

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Fetch account by ID
func GetAccount(id int) (*models.Account, error) {
	var acc models.Account
	err := DB.QueryRow(context.Background(), "SELECT id, name, balance FROM accounts WHERE id = $1", id).
		Scan(&acc.ID, &acc.Name, &acc.Balance)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

// Update Balance function for deposits & withdrawals
func UpdateBalance(accountID int, amount float64, operation string) error {

	_, err := DB.Begin(context.Background())
	if err != nil {
		return err
	}

	var query string
	if operation == "deposit" {
		query = "UPDATE accounts SET balance = balance + $1 WHERE id = $2"
	} else if operation == "withdraw" {
		query = "UPDATE accounts SET balance = balance - $1 WHERE id = $2"
	}

	err = AddTransaction(accountID, amount, operation)

	if err != nil {
		return err
	}

	_, err = DB.Exec(context.Background(), query, amount, accountID)
	return err
}

// Add Transaction Record
func AddTransaction(accountID int, amount float64, txType string) error {
	l := fmt.Sprintf("INSERT INTO transactions (account_id, amount, type) VALUES (%v, %v, %v)\n", accountID, amount, txType)
	log.Println(l)
	_, err := DB.Exec(context.Background(), "INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, $3)", accountID, amount, txType)
	return err
}

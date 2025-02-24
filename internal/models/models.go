package models

import "time"

// Account represents a bank account
type Account struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

// Transaction represents a bank transaction
type Transaction struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"` // "account_creation", "deposit", "withdraw"
	CreatedAt time.Time `json:"created_at"`
}

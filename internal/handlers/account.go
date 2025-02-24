package handlers

import (
	"banking-ledger-service/internal/models"
	"banking-ledger-service/internal/queue"
	"banking-ledger-service/internal/storage"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

// CreateAccount API handler
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc models.Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if account already exists
	var exists bool
	err := storage.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM accounts WHERE name=$1)", acc.Name).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Account already exists", http.StatusConflict)
		return
	}

	// Prepare the transaction message as JSON string
	messageData := map[string]interface{}{
		"type":    "account_creation",
		"name":    acc.Name,
		"balance": acc.Balance,
	}

	messageBytes, err := json.Marshal(messageData) // Convert to JSON string
	if err != nil {
		http.Error(w, "Failed to serialize message", http.StatusInternalServerError)
		return
	}
	messageString := string(messageBytes) // Convert bytes to string

	// Send message to RabbitMQ
	err = queue.PublishMessage(messageString)
	if err != nil {
		http.Error(w, "Failed to queue account creation", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Account creation request sent to queue"})
}

// GetAccount API handler
func GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	accID := r.URL.Query().Get("id")
	if accID == "" {
		http.Error(w, "Account ID is required", http.StatusBadRequest)
		return
	}

	// Fetch account details
	id, err := strconv.Atoi(accID)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	// Ensure account exists
	account, err := storage.GetAccount(id)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(account)
}

package handlers

import (
	"banking-ledger-service/internal/models"
	"banking-ledger-service/internal/queue"
	"banking-ledger-service/internal/storage"
	"encoding/json"
	"net/http"
)

// Deposit API handler
func Deposit(w http.ResponseWriter, r *http.Request) {
	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if tx.Amount <= 0 { // Prevent negative deposits
		http.Error(w, "Deposit amount must be greater than zero", http.StatusBadRequest)
		return
	}

	// Ensure account exists
	_, err := storage.GetAccount(tx.AccountID)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Prepare message as JSON string
	messageData := map[string]interface{}{
		"type":       "deposit",
		"account_id": tx.AccountID,
		"amount":     tx.Amount,
	}

	// Convert to JSON string
	messageBytes, err := json.Marshal(messageData)
	if err != nil {
		http.Error(w, "Failed to serialize message", http.StatusInternalServerError)
		return
	}
	messageString := string(messageBytes)

	// Send to RabbitMQ
	err = queue.PublishMessage(messageString)
	if err != nil {
		http.Error(w, "Failed to queue deposit transaction", http.StatusInternalServerError)
		return
	}

	// Respond to client
	json.NewEncoder(w).Encode(map[string]string{"message": "Deposit request sent to queue"})
}

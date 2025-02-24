package main

import (
	"banking-ledger-service/internal/handlers"
	"banking-ledger-service/internal/queue"
	"banking-ledger-service/internal/storage"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if available
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize PostgreSQL database connection
	storage.InitDB()

	// Initialize MongoDB database connection
	storage.InitMongoDB()

	// Initialize RabbitMQ connection
	queue.InitRabbitMQ()

	// Set up HTTP handlers for account creation and transactions
	http.HandleFunc("/accounts/create", handlers.CreateAccount)
	http.HandleFunc("/accounts/balance", handlers.GetAccountBalance)
	http.HandleFunc("/transactions/deposit", handlers.Deposit)
	http.HandleFunc("/transactions/withdraw", handlers.Withdraw)

	// Start the API server on port 8080
	log.Println("API Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

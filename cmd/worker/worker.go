package main

import (
	"banking-ledger-service/internal/queue"
	"banking-ledger-service/internal/storage"
	"encoding/json"
	"log"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

// ProcessTransaction handles messages from RabbitMQ
func ProcessTransaction(msg amqp091.Delivery) {
	// Parse message body
	var data map[string]interface{}
	err := json.Unmarshal(msg.Body, &data)
	if err != nil {
		log.Println("Failed to parse message:", err)
		return
	}

	log.Println("Processing transaction:", data)

	// Check transaction type
	txType, ok := data["type"].(string)
	if !ok || txType == "" {
		log.Println("Invalid transaction type")
		return
	}

	switch txType {
	case "account_creation":
		// Create account
		name := data["name"].(string)
		balance := data["balance"].(float64)
		accountID, err := storage.CreateAccount(name, balance)
		if err != nil {
			log.Println("Account creation failed:", err)
		} else {
			log.Println("Account created successfully")
			storage.LogTransactionToMongo(accountID, balance, "account_creation")
		}
	case "deposit":
		// Deposit funds
		accountID := int(data["account_id"].(float64))
		amount := data["amount"].(float64)
		if err := storage.UpdateBalance(accountID, amount, "deposit"); err != nil {
			log.Println("Deposit failed:", err)
		} else {
			log.Println("Deposit successful")
			storage.LogTransactionToMongo(accountID, amount, "deposit")
		}
	case "withdraw":
		// Withdraw funds
		accountID := int(data["account_id"].(float64))
		amount := data["amount"].(float64)
		account, err := storage.GetAccount(accountID)
		if err != nil || account.Balance < amount {
			log.Println("Withdrawal failed: Insufficient balance")
			return
		}
		if err := storage.UpdateBalance(accountID, amount, "withdraw"); err != nil {
			log.Println("Withdrawal failed:", err)
		} else {
			log.Println("Withdrawal successful")
			storage.LogTransactionToMongo(accountID, amount, "withdraw")
		}
	default:
		// Unknown transaction type
		log.Println("Unknown transaction type:", txType)
	}

	msg.Ack(false)
}

func main() {
	// Initialize storage and queue connections

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	storage.InitDB()
	storage.InitMongoDB()
	queue.InitRabbitMQ()

	log.Println("Worker started, waiting for messages...")

	messages, err := queue.ConsumeMessages()
	if err != nil {
		log.Fatal("Failed to consume messages:", err)
	}

	for msg := range messages {
		go ProcessTransaction(msg)
	}
}

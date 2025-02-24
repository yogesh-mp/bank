package storage

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var transactionCollection *mongo.Collection

// InitMongoDB initializes MongoDB connection
func InitMongoDB() {
	mongoURI := os.Getenv("MONGO_URI")
	// If no URI is provided, construct one from environment variables
	if mongoURI == "" {
		mongoUser := os.Getenv("MONGO_USER")
		mongoPass := os.Getenv("MONGO_PASSWORD")
		mongoHost := os.Getenv("MONGO_HOST")

		if mongoHost == "" {
			mongoHost = "localhost"
		}

		if mongoUser != "" && mongoPass != "" {
			mongoURI = "mongodb://" + mongoUser + ":" + mongoPass + "@" + mongoHost + ":27017"
		} else {
			mongoURI = "mongodb://" + mongoHost + ":27017"
		}
	}
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Verify connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("MongoDB connection test failed:", err)
	}

	// Set global variables
	log.Println("Connected to MongoDB!")
	log.Println(mongoClient)
	mongoClient = client
	transactionCollection = client.Database("banking_ledger").Collection("transactions")
}

// LogTransactionToMongo stores transaction logs in MongoDB
func LogTransactionToMongo(accountID int, amount float64, txType string) {
	// Create a new document
	doc := bson.M{
		"account_id": accountID,
		"amount":     amount,
		"type":       txType,
		"timestamp":  time.Now(),
	}
	// Insert the document into the collection
	_, err := transactionCollection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Println("Failed to insert transaction log into MongoDB:", err)
	} else {
		log.Println("Transaction logged successfully in MongoDB")
	}
}

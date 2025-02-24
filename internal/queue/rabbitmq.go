package queue

import (
	"fmt"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection
var channel *amqp091.Channel
var queueName = "transactions"

// Initialize RabbitMQ connection
func InitRabbitMQ() {
	var err error
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqUser := os.Getenv("RABBITMQ_USER")
	rabbitmqPass := os.Getenv("RABBITMQ_PASSWORD")

	// Construct the connection URL using the environment variables
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:5672/", rabbitmqUser, rabbitmqPass, rabbitmqHost)

	// Connect to RabbitMQ
	conn, err = amqp091.Dial(amqpURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	// Open a channel
	channel, err = conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}

	// Declare a queue
	_, err = channel.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,
	)

	// Check for errors
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	log.Println("RabbitMQ initialized successfully!")
}

// Publish a message to RabbitMQ
func PublishMessage(message string) error {
	err := channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Println("Failed to publish message:", err)
		return err
	}

	log.Println("Message published to queue:", message)
	return nil
}

// ConsumeMessages consumes messages from the queue
func ConsumeMessages() (<-chan amqp091.Delivery, error) {
	messages, err := channel.Consume(
		queueName,
		"",
		false, // Manual acknowledgment
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to consume messages:", err)
		return nil, err
	}
	return messages, nil
}

package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	// Connect to RabbitMQ container
	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"notification-queue", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Publish a message
	// body := []byte("A new message")
	// err = ch.Publish(
	// 	"",     // exchange
	// 	q.Name, // routing key
	// 	false,  // mandatory
	// 	false,  // immediate
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        body,
	// 	})
	// if err != nil {
	// 	log.Fatalf("Failed to publish a message: %v", err)
	// }

	// Receive a message
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Waiting for notifications")

	// Implement Notification Logic
	for msg := range msgs {
		log.Printf("Received a notification: %s", msg.Body)
	}
}

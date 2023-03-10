// Package notification implements the access to RabbitMQ with the notification queue
package notification

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

var (
	rabbitConn  *amqp.Connection
	rabbitCh    *amqp.Channel
	rabbitQueue amqp.Queue
)

// Connect establishes a connection to the RabbitMQ server and creates a channel and a queue.
func Connect() error {
	var err error
	// Connect to RabbitMQ container
	rabbitConn, err = amqp.Dial(os.Getenv("AMQP_SERVER_URL"))
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err.Error())
		return err
	}

	// Create a channel
	rabbitCh, err = rabbitConn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel into RabbitMQ:", err.Error())
		return err
	}

	// Declare the notification queue
	rabbitQueue, err = rabbitCh.QueueDeclare(
		"notification-queue", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	return err
}

// PublishNotification publishes a message to the notification queue.
func PublishNotification(msg []byte) error {
	err := rabbitCh.Publish(
		"",               // exchange
		rabbitQueue.Name, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	return err
}

package main

import (
	"log"

	"github.com/PixDale/sh-code-challenge/api"
	"github.com/PixDale/sh-code-challenge/api/notification"
)

func main() {
	err := notification.Connect()
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	api.Run()
}

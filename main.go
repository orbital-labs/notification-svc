package main

import (
	"fmt"
	"log"
	"os"

	"github.com/orbital-labs/notification-svc/rabbit"
)

func main() {
	err := initEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rmqURL := os.Getenv("RABBITMQ_URL")

	rmqMataairOpt := rabbit.RabbitMQOption{
		RabbitMQURL:  rmqURL,
		ExchangeName: "notification_mataair",
		ExchangeType: "direct",
		RoutingKey:   "notification.mataair",
		Durable:      true,
		Exclusive:    false,
	}

	rmqMataair := rabbit.NewRabbitMQConsumer(rmqMataairOpt)
	// go rmqMataair.ServeConsume(rmqMataair.ListenRequest)
	rmqMataair.ServeConsume(rmqMataair.ListenRequest)

	fmt.Println("Copycat Background Job is now running")

}

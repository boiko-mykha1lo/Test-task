package main

import (
	"log"
	"time"

	"github.com/yourname/notifications/internal/config"
	"github.com/yourname/notifications/internal/consumer"
)

func main() {
	cfg := config.Load()

	time.Sleep(10 * time.Second)

	c, err := consumer.NewConsumer(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	if err := c.Start(); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
}

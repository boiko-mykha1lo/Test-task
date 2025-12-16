package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ProductMessage struct {
	Action      string `json:"action"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name,omitempty"`
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		"product_events",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *Consumer) Start() error {
	msgs, err := c.channel.Consume(
		"product_events",
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Println("Notifications service started. Waiting for messages...")

	for msg := range msgs {
		c.handleMessage(msg.Body)
	}

	return nil
}

func (c *Consumer) handleMessage(body []byte) {
	var msg ProductMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	switch msg.Action {
	case "created":
		log.Printf("✅ Product CREATED: ID=%d, Name=%s", msg.ProductID, msg.ProductName)
	case "deleted":
		log.Printf("🗑️  Product DELETED: ID=%d", msg.ProductID)
	default:
		log.Printf("Unknown action: %s", msg.Action)
	}
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}

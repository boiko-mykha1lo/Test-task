package broker

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker interface {
	PublishProductCreated(productID int, productName string) error
	PublishProductDeleted(productID int) error
	Close() error
}

type rabbitMQBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type ProductMessage struct {
	Action      string `json:"action"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name,omitempty"`
}

func NewRabbitMQBroker(url string) (MessageBroker, error) {
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

	return &rabbitMQBroker{
		conn:    conn,
		channel: ch,
	}, nil
}

func (b *rabbitMQBroker) PublishProductCreated(productID int, productName string) error {
	msg := ProductMessage{
		Action:      "created",
		ProductID:   productID,
		ProductName: productName,
	}

	return b.publish(msg)
}

func (b *rabbitMQBroker) PublishProductDeleted(productID int) error {
	msg := ProductMessage{
		Action:    "deleted",
		ProductID: productID,
	}

	return b.publish(msg)
}

func (b *rabbitMQBroker) publish(msg ProductMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = b.channel.Publish(
		"",
		"product_events",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message: %s", string(body))
	return nil
}

func (b *rabbitMQBroker) Close() error {
	if err := b.channel.Close(); err != nil {
		return err
	}
	return b.conn.Close()
}

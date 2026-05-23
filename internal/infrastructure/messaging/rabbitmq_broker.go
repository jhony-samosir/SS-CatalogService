package messaging

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"ss-catalog-service/internal/worker"
)

const (
	exchangeName = "samstore.events"
	exchangeType = "topic"
)

// RabbitMQBroker implements worker.MessageBroker using RabbitMQ.
// It maintains a persistent connection with automatic reconnection.
type RabbitMQBroker struct {
	url        string
	conn       *amqp.Connection
	ch         *amqp.Channel
	mu         sync.Mutex
	maxRetries int
}

// NewRabbitMQBroker creates and connects a new RabbitMQ broker.
// It declares the topic exchange on startup.
func NewRabbitMQBroker(url string) (worker.MessageBroker, error) {
	b := &RabbitMQBroker{
		url:        url,
		maxRetries: 5,
	}
	if err := b.connect(); err != nil {
		return nil, fmt.Errorf("rabbitmq: initial connection failed: %w", err)
	}
	return b, nil
}

// connect (re)establishes the connection and channel, and redeclares the exchange.
func (b *RabbitMQBroker) connect() error {
	conn, err := amqp.Dial(b.url)
	if err != nil {
		return fmt.Errorf("rabbitmq: dial failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("rabbitmq: open channel failed: %w", err)
	}

	// Enable publisher confirms for reliable delivery
	if err := ch.Confirm(false); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("rabbitmq: confirm mode failed: %w", err)
	}

	// Declare the durable topic exchange (idempotent)
	if err := ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("rabbitmq: exchange declare failed: %w", err)
	}

	b.conn = conn
	b.ch = ch
	slog.Info("RabbitMQ connected", "exchange", exchangeName)
	return nil
}

// ensureConnected checks connection health and reconnects if necessary.
func (b *RabbitMQBroker) ensureConnected() error {
	if b.conn != nil && !b.conn.IsClosed() {
		return nil
	}

	slog.Warn("RabbitMQ connection lost, reconnecting...")
	var lastErr error
	for i := 0; i < b.maxRetries; i++ {
		backoff := time.Duration(1<<uint(i)) * time.Second
		time.Sleep(backoff)
		if err := b.connect(); err != nil {
			lastErr = err
			slog.Warn("RabbitMQ reconnect attempt failed", "attempt", i+1, "error", err)
			continue
		}
		return nil
	}
	return fmt.Errorf("rabbitmq: failed to reconnect after %d attempts: %w", b.maxRetries, lastErr)
}

// Publish sends an event to the samstore.events exchange.
// The routing key is derived from the eventType: "catalog.productcreated" → "catalog.product.created"
func (b *RabbitMQBroker) Publish(ctx context.Context, eventType string, payload []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err := b.ensureConnected(); err != nil {
		return err
	}

	routingKey := buildRoutingKey("catalog", eventType)

	confirms := b.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	err := b.ch.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         payload,
			Headers: amqp.Table{
				"event_type": eventType,
				"service":    "ss-catalog-service",
			},
		},
	)
	if err != nil {
		// Mark connection as dead so next call reconnects
		b.conn.Close()
		return fmt.Errorf("rabbitmq: publish failed: %w", err)
	}

	// Wait for broker acknowledgement (publisher confirms)
	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			return fmt.Errorf("rabbitmq: broker nack'd message (routing_key=%s)", routingKey)
		}
	case <-ctx.Done():
		return fmt.Errorf("rabbitmq: context cancelled while waiting for confirm")
	case <-time.After(5 * time.Second):
		return fmt.Errorf("rabbitmq: publish confirm timeout (routing_key=%s)", routingKey)
	}

	slog.Debug("Event published", "routing_key", routingKey, "bytes", len(payload))
	return nil
}

// Close gracefully closes the channel and connection.
func (b *RabbitMQBroker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.ch != nil {
		b.ch.Close()
	}
	if b.conn != nil && !b.conn.IsClosed() {
		b.conn.Close()
	}
	slog.Info("RabbitMQ connection closed")
}

// buildRoutingKey converts "ProductCreated" → "catalog.product.created"
func buildRoutingKey(service, eventType string) string {
	// Insert dots before uppercase letters: "ProductCreated" → "product.created"
	var result strings.Builder
	for i, r := range eventType {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('.')
		}
		result.WriteRune(r)
	}
	return service + "." + strings.ToLower(result.String())
}

package messaging

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// InboxEventModel represents a received message to ensure idempotency.
type InboxEventModel struct {
	MessageID     string    `gorm:"primaryKey;type:varchar(255)"`
	EventType     string    `gorm:"type:varchar(255);not null"`
	AggregateType string    `gorm:"type:varchar(100)"`
	Payload       []byte    `gorm:"type:jsonb;not null"`
	ProcessedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Status        string    `gorm:"type:varchar(50);not null;default:'processed'"`
}

func (InboxEventModel) TableName() string {
	return "inbox_events"
}

// InboxConsumer handles consuming messages from RabbitMQ.
type InboxConsumer struct {
	url      string
	queue    string
	db       *gorm.DB
	conn     *amqp.Connection
	ch       *amqp.Channel
}

// NewInboxConsumer creates a new RabbitMQ consumer for the inbox pattern.
func NewInboxConsumer(url string, db *gorm.DB) *InboxConsumer {
	return &InboxConsumer{
		url:   url,
		queue: "catalog-service.cart-events",
		db:    db,
	}
}

// Start connects to RabbitMQ and starts consuming messages.
// It blocks until the context is canceled.
func (c *InboxConsumer) Start(ctx context.Context) error {
	for {
		err := c.consume(ctx)
		if err != nil {
			slog.Error("Inbox consumer error", "error", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Second): // Backoff before reconnecting
			slog.Info("Reconnecting inbox consumer...")
		}
	}
}

func (c *InboxConsumer) consume(ctx context.Context) error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	defer conn.Close()
	c.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	c.ch = ch

	// Ensure exchange exists (should be created by publisher, but just in case)
	err = ch.ExchangeDeclare(
		"samstore.events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		c.queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// Bind queue to exchange with routing key cart.*
	err = ch.QueueBind(
		q.Name,
		"cart.*",
		"samstore.events",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer tag
		false, // auto-ack (we use manual ack)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	slog.Info("Inbox consumer started", "queue", q.Name)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Inbox consumer shutting down")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return nil // channel closed
			}
			c.processMessage(ctx, msg)
		}
	}
}

func (c *InboxConsumer) processMessage(ctx context.Context, msg amqp.Delivery) {
	messageID := msg.MessageId
	if messageID == "" {
		// Fallback if publisher didn't provide message ID
		// In a real system, you might want to reject this or generate a hash of the body.
		slog.Warn("Received message without MessageId", "routing_key", msg.RoutingKey)
		// We'll nack and discard it for safety, or we could generate one.
		// For simplicity here, let's reject it to avoid processing duplicates silently if it's missing.
		// However, amqplib in Node.js might not send messageId automatically unless specified.
		// Let's use correlationId or a hash of the body as a fallback.
		// Assuming publisher will set messageId.
		msg.Reject(false)
		return
	}

	eventType := ""
	if et, ok := msg.Headers["event_type"].(string); ok {
		eventType = et
	} else {
		eventType = msg.RoutingKey
	}

	slog.Info("Processing incoming event", "message_id", messageID, "routing_key", msg.RoutingKey)

	// Idempotency Check & Processing within a transaction
	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if already processed
		var count int64
		if err := tx.Model(&InboxEventModel{}).Where("message_id = ?", messageID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			// Already processed, just ack
			slog.Debug("Message already processed (idempotent)", "message_id", messageID)
			return nil
		}

		// Insert into inbox to lock this message ID
		inboxEvent := &InboxEventModel{
			MessageID: messageID,
			EventType: eventType,
			Payload:   msg.Body,
			Status:    "processed",
		}
		
		// Use Clauses to handle concurrent inserts gracefully (ON CONFLICT DO NOTHING)
		err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(inboxEvent).Error
		if err != nil {
			return err
		}
		
		// Ensure it was actually inserted, if rows affected is 0, another concurrent worker might have inserted it
		if tx.RowsAffected == 0 {
			slog.Debug("Message already processed concurrently (idempotent)", "message_id", messageID)
			return nil
		}

		// --- Process Domain Logic Here ---
		// e.g. Handle cart.checked_out to validate or log stock
		if eventType == "cart.checked_out" {
			var payload struct {
				CartID int `json:"cart_id"`
			}
			if err := json.Unmarshal(msg.Body, &payload); err == nil {
				slog.Info("Cart checked out, trigger inventory check", "cart_id", payload.CartID)
				// Call Inventory Usecase...
			}
		}

		return nil
	})

	if err != nil {
		slog.Error("Failed to process message", "message_id", messageID, "error", err)
		// Nack and requeue
		msg.Nack(false, true)
	} else {
		// Ack
		msg.Ack(false)
	}
}

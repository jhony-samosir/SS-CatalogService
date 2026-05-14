package domain

import (
	"context"
	"encoding/json"
	"time"
)

// AuditLog represents a central audit trail for entity mutations.
type AuditLog struct {
	BaseEntity
	EntityType  string
	EntityID    int
	Action      string
	OldData     json.RawMessage
	NewData     json.RawMessage
	PerformedBy string
	PerformedAt time.Time
	IPAddress   string
	UserAgent   string
}

// OutboxStatus defines the state of an outbox event.
type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "pending"
	OutboxStatusPublished OutboxStatus = "published"
	OutboxStatusFailed    OutboxStatus = "failed"
)

// OutboxEvent represents a transactional outbox record.
type OutboxEvent struct {
	BaseEntity
	EventType     string
	AggregateType string
	AggregateID   int
	Payload       json.RawMessage
	Status        OutboxStatus
	RetryCount    int
	PublishedAt   *time.Time
	ErrorMessage  string
}

// OutboxRepository defines the contract for outbox event management.
type OutboxRepository interface {
	Save(ctx context.Context, event *OutboxEvent) error
	FetchPending(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkAsPublished(ctx context.Context, id int) error
	MarkAsFailed(ctx context.Context, id int, errorMessage string) error
}

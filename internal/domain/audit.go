package domain

import (
	"encoding/json"
	"time"
)

// Seller represents the vendor information.
type Seller struct {
	BaseEntity
	Name       string
	Code       string
	IsActive   bool
	VerifiedAt *time.Time
}

// SellerProduct maps which sellers are authorized to sell each product.
type SellerProduct struct {
	BaseEntity
	SellerID   int
	ProductID  int
	IsActive   bool
	ApprovedAt *time.Time
	ApprovedBy string
}

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

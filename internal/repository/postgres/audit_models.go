package postgres

import (
	"encoding/json"
	"time"

	"ss-catalog-service/internal/domain"

	"gorm.io/datatypes"
)

// AuditLogModel represents the database schema for audit_logs.
type AuditLogModel struct {
	BaseModel
	EntityType  string         `gorm:"type:varchar(100);not null;index"`
	EntityID    int            `gorm:"not null;index"`
	Action      string         `gorm:"type:varchar(50);not null"`
	OldData     datatypes.JSON `gorm:"type:jsonb"`
	NewData     datatypes.JSON `gorm:"type:jsonb"`
	PerformedBy string         `gorm:"type:varchar(255);index"`
	PerformedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;index"`
	IPAddress   string         `gorm:"type:inet"`
	UserAgent   string         `gorm:"type:text"`
}

func (AuditLogModel) TableName() string { return "audit_logs" }

// OutboxEventModel represents the database schema for outbox_events.
type OutboxEventModel struct {
	BaseModel
	EventType     string         `gorm:"type:varchar(255);not null"`
	AggregateType string         `gorm:"type:varchar(100);not null;index"`
	AggregateID   int            `gorm:"not null;index"`
	Payload       datatypes.JSON `gorm:"type:jsonb;not null"`
	Status        string         `gorm:"type:varchar(50);not null;default:'pending';index"`
	RetryCount    int            `gorm:"not null;default:0"`
	PublishedAt   *time.Time     `gorm:"null"`
	ErrorMessage  string         `gorm:"type:text"`
}

func (OutboxEventModel) TableName() string { return "outbox_events" }

func (m *OutboxEventModel) ToDomain() domain.OutboxEvent {
	return domain.OutboxEvent{
		BaseEntity: domain.BaseEntity{
			ID:        m.ID,
			PublicID:  m.PublicID,
			CreatedAt: m.CreatedAt,
			CreatedBy: m.CreatedBy,
			UpdatedAt: m.UpdatedAt,
			UpdatedBy: m.UpdatedBy,
		},
		EventType:     m.EventType,
		AggregateType: m.AggregateType,
		AggregateID:   m.AggregateID,
		Payload:       json.RawMessage(m.Payload),
		Status:        domain.OutboxStatus(m.Status),
		RetryCount:    m.RetryCount,
		PublishedAt:   m.PublishedAt,
		ErrorMessage:  m.ErrorMessage,
	}
}

func FromOutboxDomain(e *domain.OutboxEvent) *OutboxEventModel {
	return &OutboxEventModel{
		BaseModel: BaseModel{
			ID:        e.ID,
			PublicID:  e.PublicID,
			CreatedAt: e.CreatedAt,
			CreatedBy: e.CreatedBy,
			UpdatedAt: e.UpdatedAt,
			UpdatedBy: e.UpdatedBy,
		},
		EventType:     e.EventType,
		AggregateType: e.AggregateType,
		AggregateID:   e.AggregateID,
		Payload:       datatypes.JSON(e.Payload),
		Status:        string(e.Status),
		RetryCount:    e.RetryCount,
		PublishedAt:   e.PublishedAt,
		ErrorMessage:  e.ErrorMessage,
	}
}

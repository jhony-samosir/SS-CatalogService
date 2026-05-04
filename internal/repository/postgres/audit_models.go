package postgres

import (
	"time"

	"gorm.io/datatypes"
)

// SellerModel represents the database schema for sellers.
type SellerModel struct {
	BaseModel
	Name       string     `gorm:"type:varchar(255);not null"`
	Code       string     `gorm:"type:varchar(100);not null;uniqueIndex"`
	IsActive   bool       `gorm:"not null;default:true"`
	VerifiedAt *time.Time `gorm:"null"`
}

func (SellerModel) TableName() string { return "sellers" }

// SellerProductModel represents the database schema for seller_products.
type SellerProductModel struct {
	BaseModel
	SellerID   int        `gorm:"not null;index"`
	ProductID  int        `gorm:"not null;index"`
	IsActive   bool       `gorm:"not null;default:true"`
	ApprovedAt *time.Time `gorm:"null"`
	ApprovedBy string     `gorm:"type:varchar(255)"`
}

func (SellerProductModel) TableName() string { return "seller_products" }

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

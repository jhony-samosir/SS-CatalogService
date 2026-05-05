package postgres

import (
	"context"
	"time"

	"ss-catalog-service/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type outboxRepository struct {
	db *gorm.DB
}

// NewOutboxRepository creates a new instance of PostgreSQL outbox repository.
func NewOutboxRepository(db *gorm.DB) domain.OutboxRepository {
	return &outboxRepository{db: db}
}

func (r *outboxRepository) Save(ctx context.Context, e *domain.OutboxEvent) error {
	model := FromOutboxDomain(e)
	db := getDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		return err
	}
	e.ID = model.ID
	e.CreatedAt = model.CreatedAt
	return nil
}

func (r *outboxRepository) FetchPending(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	var models []OutboxEventModel
	db := getDB(ctx, r.db)

	// Best Practice: Use FOR UPDATE SKIP LOCKED to prevent multiple workers from picking up the same events.
	err := db.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("status = ? AND retry_count < ?", string(domain.OutboxStatusPending), 5).
		Limit(limit).
		Order("created_at asc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	events := make([]domain.OutboxEvent, len(models))
	for i, m := range models {
		events[i] = m.ToDomain()
	}
	return events, nil
}

func (r *outboxRepository) MarkAsPublished(ctx context.Context, id int) error {
	db := getDB(ctx, r.db)
	now := time.Now()

	return db.Model(&OutboxEventModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       string(domain.OutboxStatusPublished),
			"published_at": &now,
		}).Error
}

func (r *outboxRepository) MarkAsFailed(ctx context.Context, id int, errorMessage string) error {
	db := getDB(ctx, r.db)

	// Logic: Increment retry_count and only mark as FAILED if max retries exceeded.
	// Otherwise, keep as PENDING to allow automatic retry by FetchPending.
	return db.Transaction(func(tx *gorm.DB) error {
		var model OutboxEventModel
		if err := tx.First(&model, id).Error; err != nil {
			return err
		}

		newStatus := string(domain.OutboxStatusPending)
		if model.RetryCount >= 4 { // Next retry will be 5
			newStatus = string(domain.OutboxStatusFailed)
		}

		return tx.Model(&model).Updates(map[string]interface{}{
			"status":        newStatus,
			"error_message": errorMessage,
			"retry_count":   gorm.Expr("retry_count + 1"),
		}).Error
	})
}

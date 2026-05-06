package postgres

import (
	"context"
	"ss-catalog-service/internal/domain"

	"gorm.io/gorm"
)

type transactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new GORM-based transaction manager.
func NewTransactionManager(db *gorm.DB) domain.TransactionManager {
	return &transactionManager{db: db}
}

type txKey struct{}

var txValueKey = txKey{}

// WithTransaction executes the given function within a database transaction.
func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Inject the transaction object into context using a custom key to avoid collisions.
		txCtx := context.WithValue(ctx, txValueKey, tx)
		return fn(txCtx)
	})
}

// Helper to get DB from context (handles both transaction and regular connection)
func getDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txValueKey).(*gorm.DB); ok {
		return tx
	}
	return defaultDB.WithContext(ctx)
}

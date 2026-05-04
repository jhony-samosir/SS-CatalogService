package domain

import "context"

// TransactionManager defines the contract for managing database transactions.
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

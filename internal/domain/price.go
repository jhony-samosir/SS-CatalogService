package domain

import (
	"context"
	"time"
)

type PriceHistory struct {
	ID        int
	ProductID int
	VariantID *int
	Price     float64
	Currency  string
	Reason    string
	CreatedAt time.Time
}

type PriceHistoryRepository interface {
	LogPriceChange(ctx context.Context, history *PriceHistory) error
	GetPriceHistory(ctx context.Context, productID int, variantID *int) ([]PriceHistory, error)
}

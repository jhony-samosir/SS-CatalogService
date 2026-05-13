package postgres

import (
	"ss-catalog-service/internal/domain"
	"time"
)

type PriceHistoryModel struct {
	ID        int       `gorm:"primaryKey"`
	ProductID int       `gorm:"index:idx_price_history_product"`
	VariantID *int      `gorm:"index:idx_price_history_variant"`
	Price     float64   `gorm:"type:decimal(10,2);not null"`
	Currency  string    `gorm:"type:varchar(3);default:'IDR'"`
	Reason    string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (PriceHistoryModel) TableName() string {
	return "product_price_history"
}

func (m *PriceHistoryModel) ToDomain() domain.PriceHistory {
	return domain.PriceHistory{
		ID:        m.ID,
		ProductID: m.ProductID,
		VariantID: m.VariantID,
		Price:     m.Price,
		Currency:  m.Currency,
		Reason:    m.Reason,
		CreatedAt: m.CreatedAt,
	}
}

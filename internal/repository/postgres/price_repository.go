package postgres

import (
	"context"
	"ss-catalog-service/internal/domain"

	"gorm.io/gorm"
)

type priceHistoryRepository struct {
	db *gorm.DB
}

func NewPriceHistoryRepository(db *gorm.DB) domain.PriceHistoryRepository {
	return &priceHistoryRepository{db: db}
}

func (r *priceHistoryRepository) LogPriceChange(ctx context.Context, history *domain.PriceHistory) error {
	model := &PriceHistoryModel{
		ProductID: history.ProductID,
		VariantID: history.VariantID,
		Price:     history.Price,
		Currency:  history.Currency,
		Reason:    history.Reason,
	}
	db := getDB(ctx, r.db)
	return db.Create(model).Error
}

func (r *priceHistoryRepository) GetPriceHistory(ctx context.Context, productID int, variantID *int) ([]domain.PriceHistory, error) {
	var models []PriceHistoryModel
	db := getDB(ctx, r.db)

	query := db.Where("product_id = ?", productID)
	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	}

	if err := query.Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	histories := make([]domain.PriceHistory, len(models))
	for i, m := range models {
		histories[i] = m.ToDomain()
	}
	return histories, nil
}

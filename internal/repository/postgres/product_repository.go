package postgres

import (
	"context"
	"errors"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of PostgreSQL product repository.
func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.Product, error) {
	var models []ProductModel
	db := getDB(ctx, r.db)

	query := db
	if p.Limit > 0 {
		query = query.Limit(p.Limit).Offset(p.Offset)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	products := make([]domain.Product, len(models))
	for i, m := range models {
		products[i] = m.ToDomain()
	}
	return products, nil
}

func (r *productRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Product, error) {
	var model ProductModel
	db := getDB(ctx, r.db)

	if err := db.Where("public_id = ?", publicID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	product := model.ToDomain()
	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, p *domain.Product) error {
	model := FromProductDomain(p) // Updated mapper name to FromProductDomain in previous turn
	db := getDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		return err
	}
	// Update domain entity with generated values (like auto-increment ID)
	p.ID = model.ID
	p.CreatedAt = model.CreatedAt
	return nil
}

package postgres

import (
	"context"
	"log/slog"
	"ss-catalog-service/internal/domain"
	"strings"

	"gorm.io/gorm"
)

type variantRepository struct {
	db *gorm.DB
}

// NewVariantRepository creates a new instance of PostgreSQL variant repository.
func NewVariantRepository(db *gorm.DB) domain.VariantRepository {
	return &variantRepository{db: db}
}

func (r *variantRepository) CreateVariant(ctx context.Context, v *domain.ProductVariant) error {
	model := FromVariantDomain(v)
	db := getDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "create_variant", "error", err)
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "sku") {
			return domain.ErrDuplicateSKU
		}
		return err
	}
	v.ID = model.ID
	v.CreatedAt = model.CreatedAt
	slog.InfoContext(ctx, "db_audit", "operation", "create_variant", "variant_id", v.ID, "sku", v.SKU)
	return nil
}

func (r *variantRepository) CreateVariantAttributes(ctx context.Context, attrs []domain.ProductVariantAttribute) error {
	if len(attrs) == 0 {
		return nil
	}
	models := make([]ProductVariantAttributeModel, len(attrs))
	for i, a := range attrs {
		models[i] = *FromVariantAttributeDomain(&a)
	}

	db := getDB(ctx, r.db)
	if err := db.Create(&models).Error; err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "create_variant_attributes", "error", err)
		return err
	}
	slog.InfoContext(ctx, "db_audit", "operation", "create_variant_attributes", "count", len(models))
	return nil
}

func (r *variantRepository) CreateVariantImages(ctx context.Context, images []domain.ProductImage) error {
	if len(images) == 0 {
		return nil
	}
	models := make([]ProductImageModel, len(images))
	for i, img := range images {
		models[i] = *FromImageDomain(&img)
	}

	db := getDB(ctx, r.db)
	if err := db.Create(&models).Error; err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "create_variant_images", "error", err)
		return err
	}
	slog.InfoContext(ctx, "db_audit", "operation", "create_variant_images", "count", len(models))
	return nil
}

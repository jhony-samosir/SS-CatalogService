package postgres

import (
	"context"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type bundleRepository struct {
	db *gorm.DB
}

func NewBundleRepository(db *gorm.DB) domain.BundleRepository {
	return &bundleRepository{db: db}
}

func (r *bundleRepository) Create(ctx context.Context, bundle *domain.ProductBundle) error {
	model := FromBundleDomain(bundle)
	db := getDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		return mapDBError(err)
	}

	bundle.ID = model.ID
	bundle.CreatedAt = model.CreatedAt
	return nil
}

func (r *bundleRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.ProductBundle, error) {
	var models []ProductBundleModel
	db := getDB(ctx, r.db)

	query := db.Preload("Items")
	if p.Limit > 0 {
		query = query.Limit(p.Limit).Offset(p.Offset)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	bundles := make([]domain.ProductBundle, len(models))
	for i, m := range models {
		bundles[i] = m.ToDomain()
	}
	return bundles, nil
}

func (r *bundleRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.ProductBundle, error) {
	var model ProductBundleModel
	db := getDB(ctx, r.db)

	if err := db.Preload("Items").Where("public_id = ?", publicID).First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	bundle := model.ToDomain()
	return &bundle, nil
}

func (r *bundleRepository) Update(ctx context.Context, bundle *domain.ProductBundle) error {
	model := FromBundleDomain(bundle)
	db := getDB(ctx, r.db)

	return db.Transaction(func(tx *gorm.DB) error {
		// Update base bundle
		if err := tx.Model(&ProductBundleModel{}).Where("id = ?", model.ID).Updates(model).Error; err != nil {
			return err
		}

		// Delete old items and insert new ones (simplest way to handle updates for many-to-many or items)
		if err := tx.Where("bundle_id = ?", model.ID).Delete(&BundleItemModel{}).Error; err != nil {
			return err
		}

		for i := range model.Items {
			model.Items[i].BundleID = model.ID
		}
		
		if len(model.Items) > 0 {
			if err := tx.Create(&model.Items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *bundleRepository) Delete(ctx context.Context, id int) error {
	db := getDB(ctx, r.db)
	return db.Delete(&ProductBundleModel{}, id).Error
}

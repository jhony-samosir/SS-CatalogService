package postgres

import (
	"context"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type brandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) domain.BrandRepository {
	return &brandRepository{db: db}
}

func (r *brandRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.Brand, error) {
	var models []BrandModel
	db := getDB(ctx, r.db)

	query := db.Model(&BrandModel{})
	if p.Limit > 0 {
		query = query.Limit(p.Limit).Offset(p.Offset)
	}

	if err := query.Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	brands := make([]domain.Brand, len(models))
	for i, m := range models {
		brands[i] = domain.Brand{
			BaseEntity: domain.BaseEntity{
				ID:        m.ID,
				PublicID:  m.PublicID,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
			Name:        m.Name,
			Slug:        m.Slug,
			LogoURL:     m.LogoURL,
			WebsiteURL:  m.WebsiteURL,
			Description: m.Description,
			IsActive:    m.IsActive,
		}
	}
	return brands, nil
}

func (r *brandRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Brand, error) {
	var m BrandModel
	db := getDB(ctx, r.db)

	if err := db.Where("public_id = ?", publicID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Brand{
		BaseEntity: domain.BaseEntity{
			ID:        m.ID,
			PublicID:  m.PublicID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		Name:        m.Name,
		Slug:        m.Slug,
		LogoURL:     m.LogoURL,
		WebsiteURL:  m.WebsiteURL,
		Description: m.Description,
		IsActive:    m.IsActive,
	}, nil
}

func (r *brandRepository) Create(ctx context.Context, brand *domain.Brand) error {
	model := &BrandModel{
		BaseModel: BaseModel{
			PublicID: brand.PublicID,
		},
		Name:        brand.Name,
		Slug:        brand.Slug,
		LogoURL:     brand.LogoURL,
		WebsiteURL:  brand.WebsiteURL,
		Description: brand.Description,
		IsActive:    brand.IsActive,
	}

	db := getDB(ctx, r.db)
	if err := db.Create(model).Error; err != nil {
		return err
	}

	brand.ID = model.ID
	brand.CreatedAt = model.CreatedAt
	return nil
}

func (r *brandRepository) Update(ctx context.Context, brand *domain.Brand) error {
	db := getDB(ctx, r.db)
	return db.Model(&BrandModel{}).
		Where("public_id = ?", brand.PublicID).
		Updates(map[string]interface{}{
			"name":        brand.Name,
			"slug":        brand.Slug,
			"logo_url":    brand.LogoURL,
			"website_url": brand.WebsiteURL,
			"description": brand.Description,
			"is_active":   brand.IsActive,
		}).Error
}

func (r *brandRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	db := getDB(ctx, r.db)
	return db.Where("public_id = ?", publicID).Delete(&BrandModel{}).Error
}

func (r *brandRepository) CountProducts(ctx context.Context, brandID int) (int64, error) {
	var count int64
	db := getDB(ctx, r.db)
	if err := db.Model(&ProductModel{}).Where("brand_id = ?", brandID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

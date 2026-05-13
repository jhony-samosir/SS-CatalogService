package postgres

import (
	"context"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new instance of PostgreSQL category repository.
func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.Category, error) {
	var models []CategoryModel
	db := getDB(ctx, r.db)

	query := db.Model(&CategoryModel{})
	if p.Limit > 0 {
		query = query.Limit(p.Limit).Offset(p.Offset)
	}

	if err := query.Order("sort_order ASC, name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]domain.Category, len(models))
	for i, m := range models {
		categories[i] = domain.Category{
			BaseEntity: domain.BaseEntity{
				ID:        m.ID,
				PublicID:  m.PublicID,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
			ParentID:    m.ParentID,
			Name:        m.Name,
			Slug:        m.Slug,
			IconURL:     m.IconURL,
			Description: m.Description,
			Level:       m.Level,
			SortOrder:   m.SortOrder,
			IsActive:    m.IsActive,
		}
	}
	return categories, nil
}

func (r *categoryRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Category, error) {
	var m CategoryModel
	db := getDB(ctx, r.db)

	if err := db.Where("public_id = ?", publicID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Category{
		BaseEntity: domain.BaseEntity{
			ID:        m.ID,
			PublicID:  m.PublicID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		ParentID:    m.ParentID,
		Name:        m.Name,
		Slug:        m.Slug,
		IconURL:     m.IconURL,
		Description: m.Description,
		Level:       m.Level,
		SortOrder:   m.SortOrder,
		IsActive:    m.IsActive,
	}, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	model := &CategoryModel{
		BaseModel: BaseModel{
			PublicID: category.PublicID,
		},
		ParentID:    category.ParentID,
		Name:        category.Name,
		Slug:        category.Slug,
		IconURL:     category.IconURL,
		Description: category.Description,
		Level:       category.Level,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
	}

	db := getDB(ctx, r.db)
	if err := db.Create(model).Error; err != nil {
		return err
	}

	category.ID = model.ID
	category.CreatedAt = model.CreatedAt
	return nil
}

package domain

import (
	"context"

	"github.com/google/uuid"
)

// Category represents the product categorization hierarchy.
type Category struct {
	BaseEntity
	ParentID    *int   `json:"parent_id,omitempty"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	IconURL     string `json:"icon_url,omitempty"`
	Description string `json:"description,omitempty"`
	Level       int    `json:"level"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`

	// Relations
	Children []Category `json:"children,omitempty"`
}

// CategoryTranslation represents localized category names.
type CategoryTranslation struct {
	BaseEntity
	CategoryID int    `json:"category_id"`
	LangCode   string `json:"lang_code"`
	Name       string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CategorySEO represents category-specific SEO metadata.
type CategorySEO struct {
	BaseEntity
	CategoryID      int    `json:"category_id"`
	LangCode        string `json:"lang_code"`
	Slug            string `json:"slug"`
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
}

// CategoryRepository defines the contract for category data access.
type CategoryRepository interface {
	FindAll(ctx context.Context, p Pagination) ([]Category, error)
	FindByID(ctx context.Context, id int) (*Category, error)
	FindByPublicID(ctx context.Context, publicID uuid.UUID) (*Category, error)
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, publicID uuid.UUID) error
	CountChildren(ctx context.Context, parentID int) (int64, error)
	CountProducts(ctx context.Context, categoryID int) (int64, error)
}

// CategoryUsecase defines the business logic for categories.
type CategoryUsecase interface {
	GetCategories(ctx context.Context, p Pagination) ([]Category, error)
	GetCategoryByPublicID(ctx context.Context, publicID uuid.UUID) (*Category, error)
	CreateCategory(ctx context.Context, category *Category) error
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, publicID uuid.UUID) error
}

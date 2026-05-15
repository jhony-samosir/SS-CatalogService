package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// BaseEntity contains common fields for enterprise-grade audit and soft-deletion.
type BaseEntity struct {
	ID        int        `json:"-"`
	PublicID  uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy string     `json:"deleted_by,omitempty"`
}

// ProductStatus defines the lifecycle state of a product.
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusActive    ProductStatus = "active"
	ProductStatusArchived  ProductStatus = "archived"
	ProductStatusSuspended ProductStatus = "suspended"
)

// Product represents the SPU-level product records.
type Product struct {
	BaseEntity
	BrandID      *int
	SellerID     *int
	Name         string
	Slug         string
	Description  string
	ShortDesc    string
	Status       ProductStatus
	PublishAt    *time.Time
	UnpublishAt  *time.Time
	IsFeatured   bool
	WeightGram   *int
	IsDigital    bool
	ImageURL     string

	// Aggregated Data
	Translation *ProductTranslation
	SEO         *ProductSEO
	Categories  []Category
	Tags        []Tag
}

// ProductRepository defines the contract for product data access.
type ProductRepository interface {
	FindAll(ctx context.Context, p Pagination) ([]Product, int64, error)
	FindByID(ctx context.Context, id int) (*Product, error)
	FindByPublicID(ctx context.Context, publicID uuid.UUID) (*Product, error)
	GetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string) (*Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Search(ctx context.Context, q GetProductSearchQuery) (*ProductSearchResult, error)
}

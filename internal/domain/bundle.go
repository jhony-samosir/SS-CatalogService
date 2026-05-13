package domain

import (
	"context"
	"github.com/google/uuid"
)

type ProductBundle struct {
	BaseEntity
	Name          string
	Slug          string
	Description   string
	PriceOverride *float64
	IsActive      bool
	Items         []BundleItem
}

type BundleItem struct {
	BundleID   int
	ProductID  *int
	VariantID  *int
	Quantity   int
	IsOptional bool
}

// --- Interfaces ---

type BundleRepository interface {
	Create(ctx context.Context, bundle *ProductBundle) error
	FindAll(ctx context.Context, p Pagination) ([]ProductBundle, error)
	FindByPublicID(ctx context.Context, publicID uuid.UUID) (*ProductBundle, error)
	Update(ctx context.Context, bundle *ProductBundle) error
	Delete(ctx context.Context, id int) error
}

type BundleUsecase interface {
	CreateBundle(ctx context.Context, bundle *ProductBundle) error
	GetBundles(ctx context.Context, p Pagination) ([]ProductBundle, error)
	GetBundleByPublicID(ctx context.Context, publicID uuid.UUID) (*ProductBundle, error)
	UpdateBundle(ctx context.Context, bundle *ProductBundle) error
}

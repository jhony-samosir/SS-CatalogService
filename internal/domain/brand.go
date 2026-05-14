package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Brand represents the product brand information.
type Brand struct {
	BaseEntity
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	LogoURL     string `json:"logo_url,omitempty"`
	WebsiteURL  string `json:"website_url,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// BrandRepository defines the contract for brand data access.
type BrandRepository interface {
	FindAll(ctx context.Context, p Pagination) ([]Brand, error)
	FindByPublicID(ctx context.Context, publicID uuid.UUID) (*Brand, error)
	Create(ctx context.Context, brand *Brand) error
	Update(ctx context.Context, brand *Brand) error
	Delete(ctx context.Context, publicID uuid.UUID) error
	CountProducts(ctx context.Context, brandID int) (int64, error)
}

// BrandUsecase defines the business logic for brands.
type BrandUsecase interface {
	GetBrands(ctx context.Context, p Pagination) ([]Brand, error)
	GetBrandByPublicID(ctx context.Context, publicID uuid.UUID) (*Brand, error)
	CreateBrand(ctx context.Context, brand *Brand) error
	UpdateBrand(ctx context.Context, brand *Brand) error
	DeleteBrand(ctx context.Context, publicID uuid.UUID) error
}

type MasterDataCacheRepository interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	InvalidateAll(ctx context.Context) error
}

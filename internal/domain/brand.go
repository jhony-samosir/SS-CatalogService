package domain

import (
	"context"

	"github.com/google/uuid"
)

// Brand represents the product brand information.
type Brand struct {
	BaseEntity
	Name        string
	Slug        string
	LogoURL     string
	WebsiteURL  string
	Description string
	IsActive    bool
}

// BrandRepository defines the contract for brand data access.
type BrandRepository interface {
	FindAll(ctx context.Context, p Pagination) ([]Brand, error)
	FindByPublicID(ctx context.Context, publicID uuid.UUID) (*Brand, error)
	Create(ctx context.Context, brand *Brand) error
}

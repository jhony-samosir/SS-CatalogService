package domain

import (
	"context"

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
}

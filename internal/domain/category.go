package domain

import (
	"context"

	"github.com/google/uuid"
)

// Category represents hierarchical product categories.
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
}

// CategoryRepository defines the contract for category data access.
type CategoryRepository interface {
	FindAll(ctx context.Context, p Pagination) ([]Category, error)
	FindByPublicID(ctx context.Context, publicID uuid.UUID) (*Category, error)
	Create(ctx context.Context, category *Category) error
}

// CategoryUsecase defines the contract for category business logic.
type CategoryUsecase interface {
	GetCategories(ctx context.Context, p Pagination) ([]Category, error)
}

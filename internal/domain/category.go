package domain

import (
	"context"

	"github.com/google/uuid"
)

// Category represents hierarchical product categories.
type Category struct {
	BaseEntity
	ParentID    *int
	Name        string
	Slug        string
	IconURL     string
	Description string
	Level       int
	SortOrder   int
	IsActive    bool
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

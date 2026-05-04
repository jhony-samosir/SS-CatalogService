package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Sentinel Errors
var (
	ErrProductNotFound = errors.New("product not found")
)

// --- Usecase Interfaces ---

// ProductCommandUsecase defines the contract for write operations on products.
type ProductCommandUsecase interface {
	CreateProduct(ctx context.Context, payload CreateProductPayload) (*Product, error)
}

// ProductQueryUsecase defines the contract for read operations on products.
type ProductQueryUsecase interface {
	GetAllProducts(ctx context.Context, p Pagination) ([]Product, error)
	GetProductByPublicID(ctx context.Context, publicID uuid.UUID) (*Product, error)
}

// --- Payload DTOs ---

// CreateProductPayload represents the data needed to create a new product.
type CreateProductPayload struct {
	Name    string
	BrandID *int
}

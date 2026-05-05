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
	GetProductDetails(ctx context.Context, query GetProductDetailsQuery) (*ProductDetailsResponse, error)
}

// --- Payload DTOs ---

// GetProductDetailsQuery represents the input for fetching localized product details.
type GetProductDetailsQuery struct {
	PublicID uuid.UUID
	LangCode string
}

// ProductDetailsResponse represents the flattened, localized product details.
type ProductDetailsResponse struct {
	PublicID        uuid.UUID
	Name            string
	Description     string
	ShortDesc       string
	BasePrice       float64
	Status          string
	MetaTitle       string
	MetaDescription string
	Categories      []string
	Tags            []string
}

// CreateProductPayload represents the data needed to create a new product.
type CreateProductPayload struct {
	Name    string
	BrandID *int
}

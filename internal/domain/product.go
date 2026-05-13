package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Sentinel Errors
var (
	ErrProductNotFound    = errors.New("product not found")
	ErrDuplicateProduct   = errors.New("product already exists")
	ErrInternalDatabase   = errors.New("internal database error")
	ErrUnauthorized       = errors.New("unauthorized: ownership verification failed")
	ErrInvalidProductName = errors.New("product name cannot be empty")
)

// --- Usecase Interfaces ---

// ProductCommandUsecase defines the contract for write operations on products.
type ProductCommandUsecase interface {
	CreateProduct(ctx context.Context, payload CreateProductPayload) (*Product, error)
	UpdateProduct(ctx context.Context, payload UpdateProductPayload) error
}

// ProductQueryUsecase defines the contract for read operations on products.
type ProductQueryUsecase interface {
	GetAllProducts(ctx context.Context, p Pagination) ([]Product, error)
	GetProductByPublicID(ctx context.Context, publicID uuid.UUID) (*Product, error)
	GetProductDetails(ctx context.Context, query GetProductDetailsQuery) (*ProductDetailsResponse, error)
	SearchProducts(ctx context.Context, q GetProductSearchQuery) (*ProductSearchResult, error)
	FacetedSearch(ctx context.Context, q GetProductSearchQuery) (*FacetedSearchResult, error)
}

type SearchFacet struct {
	Name   string `json:"name"`
	Values []struct {
		Value string `json:"value"`
		Count int    `json:"count"`
	} `json:"values"`
}

type FacetedSearchResult struct {
	Items      []Product     `json:"items"`
	Facets     []SearchFacet `json:"facets"`
	TotalHint  int64         `json:"total_hint"`
	NextCursor *string       `json:"next_cursor"`
}

type SearchRepository interface {
	IndexProduct(ctx context.Context, product Product) error
	Search(ctx context.Context, q GetProductSearchQuery) (*FacetedSearchResult, error)
}

// ProductCacheRepository defines the contract for in-memory caching operations.
type ProductCacheRepository interface {
	GetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string) (*ProductDetailsResponse, error)
	SetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string, product *ProductDetailsResponse) error
	InvalidateProductDetails(ctx context.Context, publicID uuid.UUID) error
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

// UpdateProductPayload represents the data needed to update an existing product.
type UpdateProductPayload struct {
	PublicID    uuid.UUID
	Name        string
	Description string
	Status      ProductStatus
}

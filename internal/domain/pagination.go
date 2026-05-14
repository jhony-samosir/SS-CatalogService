package domain

import "errors"

// Pagination represents a bounded query request (offset-based).
type Pagination struct {
	Limit  int
	Offset int
}

// GetProductSearchQuery represents all optional filters for the product search API.
type GetProductSearchQuery struct {
	// Keyword for full-text search. Nil if no text search is requested.
	Keyword *string

	// CategorySlug filters products by their category slug.
	CategorySlug *string

	// BrandID filters products by their internal brand ID.
	BrandID *int

	// MinPrice is the minimum price filter (inclusive).
	MinPrice *float64

	// MaxPrice is the maximum price filter (inclusive).
	MaxPrice *float64

	// Status filters products by their lifecycle state (active, draft, etc.).
	Status *ProductStatus

	// Cursor is an opaque token for keyset-based pagination.
	Cursor *string

	// Limit is the maximum number of items to return.
	Limit int
}

// ProductSearchResult wraps the paginated search response.
type ProductSearchResult struct {
	// Items are the products found in this page.
	Items []Product

	// NextCursor is the token for the next page. Nil if last page.
	NextCursor *string

	// TotalHint is an approximate total count of matching products for UI feedback.
	TotalHint int64
}

// ErrInvalidCursor is returned when the provided pagination cursor is malformed or invalid.
var ErrInvalidCursor = errors.New("invalid or expired pagination cursor")

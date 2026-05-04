package domain

import (
	"context"
	"errors"
)

// ProductVariant represents a SKU-level sellable product unit.
type ProductVariant struct {
	BaseEntity
	ProductID  int
	SKU        string
	Barcode    string
	Name       string
	IsDefault  bool
	IsActive   bool
	WeightGram *int
	SortOrder  int
}

// Sentinel Errors
var (
	ErrDuplicateSKU = errors.New("sku already exists")
)

// --- DTOs ---

type CreateVariantPayload struct {
	ProductID  int
	SKU        string
	Barcode    string
	Name       string
	IsDefault  bool
	WeightGram *int
	Attributes []VariantAttributePayload
	Images     []VariantImagePayload
}

type VariantAttributePayload struct {
	AttributeID      int
	AttributeValueID int
}

type VariantImagePayload struct {
	URL       string
	AltText   string
	IsPrimary bool
}

// --- Interfaces ---

type VariantRepository interface {
	CreateVariant(ctx context.Context, variant *ProductVariant) error
	CreateVariantAttributes(ctx context.Context, attrs []ProductVariantAttribute) error
	CreateVariantImages(ctx context.Context, images []ProductImage) error
}

type VariantCommandUsecase interface {
	CreateProductVariant(ctx context.Context, payload CreateVariantPayload) (*ProductVariant, error)
}

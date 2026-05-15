package domain

import (
	"context"
	"github.com/google/uuid"
)

// AttributeInputType defines how the attribute value is collected.
type AttributeInputType string

const (
	InputTypeText        AttributeInputType = "text"
	InputTypeSelect      AttributeInputType = "select"
	InputTypeMultiSelect AttributeInputType = "multiselect"
	InputTypeBoolean     AttributeInputType = "boolean"
	InputTypeNumber      AttributeInputType = "number"
)

// ProductAttribute represents attribute definitions (e.g., Color, Size).
type ProductAttribute struct {
	BaseEntity
	Name      string             `json:"name"`
	Code      string             `json:"code"`
	InputType AttributeInputType `json:"input_type"`
	IsVariant bool               `json:"is_variant"`
	SortOrder int                `json:"sort_order"`
	
	// Relations
	Values []AttributeValue `json:"values,omitempty"`
}

// AttributeValue represents possible values for each attribute.
type AttributeValue struct {
	BaseEntity
	AttributeID int    `json:"attribute_id"`
	Value       string `json:"value"`
	ColorHex    string `json:"color_hex,omitempty"`
	SortOrder   int    `json:"sort_order"`
}

// Tag represents flat keyword tags for product discovery.
type Tag struct {
	BaseEntity
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ProductVariantAttribute maps a variant to its specific attribute values.
type ProductVariantAttribute struct {
	BaseEntity
	VariantID        int
	AttributeID      int
	AttributeValueID int
}

// AttributeRepository defines the contract for attribute data access.
type AttributeRepository interface {
	FindAll(ctx context.Context, p Pagination) ([]ProductAttribute, error)
	Count(ctx context.Context) (int64, error)
	FindByPublicID(ctx context.Context, publicID uuid.UUID) (*ProductAttribute, error)
	Create(ctx context.Context, attr *ProductAttribute) error
	Update(ctx context.Context, attr *ProductAttribute) error
	Delete(ctx context.Context, publicID uuid.UUID) error
	CountUsage(ctx context.Context, attrID int) (int64, error)
}

// TagRepository defines the contract for tag data access.
type TagRepository interface {
	FindAll(ctx context.Context, p Pagination) ([]Tag, int64, error)
	Create(ctx context.Context, tag *Tag) error
	Delete(ctx context.Context, publicID uuid.UUID) error
}

// AttributeUsecase defines the business logic for attributes.
type AttributeUsecase interface {
	GetAttributes(ctx context.Context, p Pagination) ([]ProductAttribute, int64, error)
	GetAttributeByPublicID(ctx context.Context, publicID uuid.UUID) (*ProductAttribute, error)
	CreateAttribute(ctx context.Context, attr *ProductAttribute) error
	UpdateAttribute(ctx context.Context, attr *ProductAttribute) error
	DeleteAttribute(ctx context.Context, publicID uuid.UUID) error
}

// TagUsecase defines the business logic for tags.
type TagUsecase interface {
	GetTags(ctx context.Context, p Pagination) ([]Tag, int64, error)
	CreateTag(ctx context.Context, tag *Tag) error
	DeleteTag(ctx context.Context, publicID uuid.UUID) error
}

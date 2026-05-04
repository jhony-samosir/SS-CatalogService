package postgres

import (
	"ss-catalog-service/internal/domain"
)

// ProductVariantModel represents the database schema for product_variants.
type ProductVariantModel struct {
	BaseModel
	ProductID  int    `gorm:"not null;index"`
	SKU        string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Barcode    string `gorm:"type:varchar(255)"`
	Name       string `gorm:"type:varchar(500)"`
	IsDefault  bool   `gorm:"not null;default:false"`
	IsActive   bool   `gorm:"not null;default:true"`
	WeightGram *int   `gorm:"null"`
	SortOrder  int    `gorm:"not null;default:0"`
}

func (ProductVariantModel) TableName() string { return "product_variants" }

// ProductAttributeModel represents the database schema for product_attributes.
type ProductAttributeModel struct {
	BaseModel
	Name      string `gorm:"type:varchar(255);not null"`
	Code      string `gorm:"type:varchar(100);not null;uniqueIndex"`
	InputType string `gorm:"type:varchar(50);not null;default:'select'"`
	IsVariant bool   `gorm:"not null;default:true"`
	SortOrder int    `gorm:"not null;default:0"`
}

func (ProductAttributeModel) TableName() string { return "product_attributes" }

// AttributeValueModel represents the database schema for attribute_values.
type AttributeValueModel struct {
	BaseModel
	AttributeID int    `gorm:"not null;index"`
	Value       string `gorm:"type:varchar(255);not null"`
	ColorHex    string `gorm:"type:varchar(7)"`
	SortOrder   int    `gorm:"not null;default:0"`
}

func (AttributeValueModel) TableName() string { return "attribute_values" }

// ProductVariantAttributeModel represents the database schema for product_variant_attributes.
type ProductVariantAttributeModel struct {
	BaseModel
	VariantID        int `gorm:"not null;index"`
	AttributeID      int `gorm:"not null"`
	AttributeValueID int `gorm:"not null"`
}

func (ProductVariantAttributeModel) TableName() string { return "product_variant_attributes" }

// ProductImageModel represents the database schema for product_images.
type ProductImageModel struct {
	BaseModel
	ProductID  *int   `gorm:"index"`
	VariantID  *int   `gorm:"index"`
	URL        string `gorm:"type:text;not null"`
	AltText    string `gorm:"type:varchar(500)"`
	SortOrder  int    `gorm:"not null;default:0"`
	IsPrimary  bool   `gorm:"not null;default:false"`
	WidthPX    int    `gorm:"null"`
	HeightPX   int    `gorm:"null"`
}

func (ProductImageModel) TableName() string { return "product_images" }

// ProductVideoModel represents the database schema for product_videos.
type ProductVideoModel struct {
	BaseModel
	ProductID   int    `gorm:"not null;index"`
	URL         string `gorm:"type:text;not null"`
	Thumbnail   string `gorm:"type:text"`
	Title       string `gorm:"type:varchar(255)"`
	DurationSec int    `gorm:"null"`
	SortOrder   int    `gorm:"not null;default:0"`
}

func (ProductVideoModel) TableName() string { return "product_videos" }

// Mapping functions can be added here as needed for each model.
func (m *ProductVariantModel) ToDomain() domain.ProductVariant {
	return domain.ProductVariant{
		BaseEntity: domain.BaseEntity{
			ID:        m.ID,
			PublicID:  m.PublicID,
			CreatedAt: m.CreatedAt,
			CreatedBy: m.CreatedBy,
			UpdatedAt: m.UpdatedAt,
			UpdatedBy: m.UpdatedBy,
			DeletedAt: &m.DeletedAt.Time,
			DeletedBy: m.DeletedBy,
		},
		ProductID:  m.ProductID,
		SKU:        m.SKU,
		Barcode:    m.Barcode,
		Name:       m.Name,
		IsDefault:  m.IsDefault,
		IsActive:   m.IsActive,
		WeightGram: m.WeightGram,
		SortOrder:  m.SortOrder,
	}
}

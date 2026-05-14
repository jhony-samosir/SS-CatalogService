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

func FromVariantDomain(v *domain.ProductVariant) *ProductVariantModel {
	return &ProductVariantModel{
		BaseModel: BaseModel{
			ID:        v.ID,
			PublicID:  v.PublicID,
			CreatedAt: v.CreatedAt,
			CreatedBy: v.CreatedBy,
			UpdatedAt: v.UpdatedAt,
			UpdatedBy: v.UpdatedBy,
			DeletedBy: v.DeletedBy,
		},
		ProductID:  v.ProductID,
		SKU:        v.SKU,
		Barcode:    v.Barcode,
		Name:       v.Name,
		IsDefault:  v.IsDefault,
		IsActive:   v.IsActive,
		WeightGram: v.WeightGram,
		SortOrder:  v.SortOrder,
	}
}

func FromVariantAttributeDomain(a *domain.ProductVariantAttribute) *ProductVariantAttributeModel {
	return &ProductVariantAttributeModel{
		BaseModel: BaseModel{
			ID:        a.ID,
			PublicID:  a.PublicID,
			CreatedAt: a.CreatedAt,
			CreatedBy: a.CreatedBy,
			UpdatedAt: a.UpdatedAt,
			UpdatedBy: a.UpdatedBy,
			DeletedBy: a.DeletedBy,
		},
		VariantID:        a.VariantID,
		AttributeID:      a.AttributeID,
		AttributeValueID: a.AttributeValueID,
	}
}

func FromImageDomain(img *domain.ProductImage) *ProductImageModel {
	return &ProductImageModel{
		BaseModel: BaseModel{
			ID:        img.ID,
			PublicID:  img.PublicID,
			CreatedAt: img.CreatedAt,
			CreatedBy: img.CreatedBy,
			UpdatedAt: img.UpdatedAt,
			UpdatedBy: img.UpdatedBy,
			DeletedBy: img.DeletedBy,
		},
		ProductID: img.ProductID,
		VariantID: img.VariantID,
		URL:       img.URL,
		AltText:   img.AltText,
		SortOrder: img.SortOrder,
		IsPrimary: img.IsPrimary,
		WidthPX:   img.WidthPX,
		HeightPX:  img.HeightPX,
	}
}

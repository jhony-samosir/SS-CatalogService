package postgres

import (
	"ss-catalog-service/internal/domain"

	"gorm.io/gorm"
)

type ProductBundleModel struct {
	BaseModel
	Name          string         `gorm:"type:varchar(255);not null"`
	Slug          string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Description   string         `gorm:"type:text"`
	PriceOverride *float64       `gorm:"type:decimal(10,2)"`
	IsActive      bool           `gorm:"default:true"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	DeletedBy     string         `gorm:"type:varchar(255)"`
	
	Items []BundleItemModel `gorm:"foreignKey:BundleID"`
}

func (ProductBundleModel) TableName() string {
	return "product_bundles"
}

type BundleItemModel struct {
	ID         int  `gorm:"primaryKey"`
	BundleID   int  `gorm:"index"`
	ProductID  *int `gorm:"index"`
	VariantID  *int `gorm:"index"`
	Quantity   int  `gorm:"default:1"`
	IsOptional bool `gorm:"default:false"`
}

func (BundleItemModel) TableName() string {
	return "product_bundle_items"
}

func (m *ProductBundleModel) ToDomain() domain.ProductBundle {
	bundle := domain.ProductBundle{
		BaseEntity: domain.BaseEntity{
			ID:        m.ID,
			PublicID:  m.PublicID,
			CreatedAt: m.CreatedAt,
			CreatedBy: m.CreatedBy,
			UpdatedAt: m.UpdatedAt,
			UpdatedBy: m.UpdatedBy,
		},
		Name:          m.Name,
		Slug:          m.Slug,
		Description:   m.Description,
		PriceOverride: m.PriceOverride,
		IsActive:      m.IsActive,
	}

	if m.DeletedAt.Valid {
		bundle.DeletedAt = &m.DeletedAt.Time
		bundle.DeletedBy = m.DeletedBy
	}

	for _, item := range m.Items {
		bundle.Items = append(bundle.Items, domain.BundleItem{
			BundleID:   item.BundleID,
			ProductID:  item.ProductID,
			VariantID:  item.VariantID,
			Quantity:   item.Quantity,
			IsOptional: item.IsOptional,
		})
	}

	return bundle
}

func FromBundleDomain(d *domain.ProductBundle) *ProductBundleModel {
	m := &ProductBundleModel{
		BaseModel: BaseModel{
			ID:        d.ID,
			PublicID:  d.PublicID,
			CreatedAt: d.CreatedAt,
			CreatedBy: d.CreatedBy,
			UpdatedBy: d.UpdatedBy,
		},
		Name:          d.Name,
		Slug:          d.Slug,
		Description:   d.Description,
		PriceOverride: d.PriceOverride,
		IsActive:      d.IsActive,
	}

	if d.UpdatedAt != nil {
		m.UpdatedAt = d.UpdatedAt
	}

	if d.DeletedAt != nil {
		m.DeletedAt = gorm.DeletedAt{Time: *d.DeletedAt, Valid: true}
		m.DeletedBy = d.DeletedBy
	}

	for _, item := range d.Items {
		m.Items = append(m.Items, BundleItemModel{
			ProductID:  item.ProductID,
			VariantID:  item.VariantID,
			Quantity:   item.Quantity,
			IsOptional: item.IsOptional,
		})
	}

	return m
}

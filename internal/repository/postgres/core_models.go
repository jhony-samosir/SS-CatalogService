package postgres

import (
	"time"

	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for GORM models.
type BaseModel struct {
	ID        int            `gorm:"primaryKey;autoIncrement"`
	PublicID  uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedBy string         `gorm:"type:varchar(255)"`
	UpdatedAt *time.Time     `gorm:"null"`
	UpdatedBy string         `gorm:"type:varchar(255)"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DeletedBy string         `gorm:"type:varchar(255)"`
}

// BrandModel represents the database schema for brands.
type BrandModel struct {
	BaseModel
	Name        string `gorm:"type:varchar(255);not null"`
	Slug        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	LogoURL     string `gorm:"type:text"`
	WebsiteURL  string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	IsActive    bool   `gorm:"not null;default:true"`
}

func (BrandModel) TableName() string { return "brands" }

// CategoryModel represents the database schema for categories.
type CategoryModel struct {
	BaseModel
	ParentID    *int   `gorm:"index"`
	Name        string `gorm:"type:varchar(255);not null"`
	Slug        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	IconURL     string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	Level       int    `gorm:"not null;default:1"`
	SortOrder   int    `gorm:"not null;default:0"`
	IsActive    bool   `gorm:"not null;default:true"`
}

func (CategoryModel) TableName() string { return "categories" }

// ProductModel represents the database schema for Products (SPU level).
type ProductModel struct {
	BaseModel
	BrandID      *int          `gorm:"index"`
	SellerID     *int          `gorm:"index"`
	Name         string        `gorm:"type:varchar(500);not null"`
	Slug         string        `gorm:"type:varchar(500);uniqueIndex;not null"`
	Description  string        `gorm:"type:text"`
	ShortDesc    string        `gorm:"type:varchar(1000)"`
	Status       string        `gorm:"type:varchar(50);not null;default:'draft'"`
	PublishAt    *time.Time    `gorm:"null"`
	UnpublishAt  *time.Time    `gorm:"null"`
	IsFeatured   bool          `gorm:"not null;default:false"`
	WeightGram   *int          `gorm:"null"`
	SearchVector string        `gorm:"type:tsvector"`

	// Associations
	Translations []ProductTranslationModel `gorm:"foreignKey:ProductID"`
	SEO          []ProductSEOModel          `gorm:"foreignKey:ProductID"`
	Categories   []CategoryModel            `gorm:"many2many:product_categories;foreignKey:ID;joinForeignKey:ProductID;References:ID;joinReferences:CategoryID"`
	Tags         []TagModel                 `gorm:"many2many:product_tags;foreignKey:ID;joinForeignKey:ProductID;References:ID;joinReferences:TagID"`

	// Optimized for single-query detail fetch
	Translation *ProductTranslationModel `gorm:"-"`
	ProductSEO  *ProductSEOModel         `gorm:"-"`
}

func (ProductModel) TableName() string { return "products" }

// ToDomain mapping functions (simplified for now)
func (m *ProductModel) ToDomain() domain.Product {
	p := domain.Product{
		BaseEntity: domain.BaseEntity{
			ID:        m.ID,
			PublicID:  m.PublicID,
			CreatedAt: m.CreatedAt,
			CreatedBy: m.CreatedBy,
			UpdatedAt: m.UpdatedAt,
			UpdatedBy: m.UpdatedBy,
			DeletedBy: m.DeletedBy,
		},
		BrandID:     m.BrandID,
		SellerID:    m.SellerID,
		Name:        m.Name,
		Slug:        m.Slug,
		Description: m.Description,
		ShortDesc:   m.ShortDesc,
		Status:      domain.ProductStatus(m.Status),
		PublishAt:   m.PublishAt,
		UnpublishAt: m.UnpublishAt,
		IsFeatured:  m.IsFeatured,
		WeightGram:  m.WeightGram,
	}

	// Map Translation
	if m.Translation != nil {
		t := m.Translation
		p.Translation = &domain.ProductTranslation{
			BaseEntity: domain.BaseEntity{
				ID:        t.ID,
				PublicID:  t.PublicID,
				CreatedAt: t.CreatedAt,
			},
			ProductID: t.ProductID,
			LangCode:  t.LangCode,
			Name:      t.Name,
			Description: t.Description,
			ShortDesc:   t.ShortDesc,
		}
	} else if len(m.Translations) > 0 {
		t := m.Translations[0]
		p.Translation = &domain.ProductTranslation{
			BaseEntity: domain.BaseEntity{
				ID:        t.ID,
				PublicID:  t.PublicID,
				CreatedAt: t.CreatedAt,
			},
			ProductID: t.ProductID,
			LangCode:  t.LangCode,
			Name:      t.Name,
			Description: t.Description,
			ShortDesc:   t.ShortDesc,
		}
	}

	// Map SEO
	if m.ProductSEO != nil {
		s := m.ProductSEO
		p.SEO = &domain.ProductSEO{
			BaseEntity: domain.BaseEntity{
				ID:        s.ID,
				PublicID:  s.PublicID,
				CreatedAt: s.CreatedAt,
			},
			ProductID:      s.ProductID,
			SEOBase: domain.SEOBase{
				LangCode:       s.LangCode,
				Slug:           s.Slug,
				MetaTitle:      s.MetaTitle,
				MetaDescription: s.MetaDescription,
			},
			CanonicalURL:   s.CanonicalURL,
			OGImageURL:     s.OGImageURL,
		}
	} else if len(m.SEO) > 0 {
		s := m.SEO[0]
		p.SEO = &domain.ProductSEO{
			BaseEntity: domain.BaseEntity{
				ID:        s.ID,
				PublicID:  s.PublicID,
				CreatedAt: s.CreatedAt,
			},
			ProductID:      s.ProductID,
			SEOBase: domain.SEOBase{
				LangCode:       s.LangCode,
				Slug:           s.Slug,
				MetaTitle:      s.MetaTitle,
				MetaDescription: s.MetaDescription,
			},
			CanonicalURL:   s.CanonicalURL,
			OGImageURL:     s.OGImageURL,
		}
	}

	// Map Categories
	if len(m.Categories) > 0 {
		p.Categories = make([]domain.Category, len(m.Categories))
		for i, c := range m.Categories {
			p.Categories[i] = domain.Category{
				BaseEntity: domain.BaseEntity{
					ID:        c.ID,
					PublicID:  c.PublicID,
					CreatedAt: c.CreatedAt,
				},
				Name:        c.Name,
				Slug:        c.Slug,
				Description: c.Description,
			}
		}
	}

	// Map Tags
	if len(m.Tags) > 0 {
		p.Tags = make([]domain.Tag, len(m.Tags))
		for i, t := range m.Tags {
			p.Tags[i] = domain.Tag{
				BaseEntity: domain.BaseEntity{
					ID:        t.ID,
					PublicID:  t.PublicID,
					CreatedAt: t.CreatedAt,
				},
				Name: t.Name,
				Slug: t.Slug,
			}
		}
	}

	return p
}

func FromProductDomain(p *domain.Product) *ProductModel {
	return &ProductModel{
		BaseModel: BaseModel{
			ID:        p.ID,
			PublicID:  p.PublicID,
			CreatedAt: p.CreatedAt,
			CreatedBy: p.CreatedBy,
			UpdatedAt: p.UpdatedAt,
			UpdatedBy: p.UpdatedBy,
			DeletedBy: p.DeletedBy,
		},
		BrandID:     p.BrandID,
		SellerID:    p.SellerID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		ShortDesc:   p.ShortDesc,
		Status:      string(p.Status),
		PublishAt:   p.PublishAt,
		UnpublishAt: p.UnpublishAt,
		IsFeatured:  p.IsFeatured,
		WeightGram:  p.WeightGram,
	}
}

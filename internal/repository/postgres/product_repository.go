package postgres

import (
	"context"
	"errors"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of PostgreSQL product repository.
func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.Product, error) {
	var models []ProductModel
	db := getDB(ctx, r.db)

	query := db
	if p.Limit > 0 {
		query = query.Limit(p.Limit).Offset(p.Offset)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	products := make([]domain.Product, len(models))
	for i, m := range models {
		products[i] = m.ToDomain()
	}
	return products, nil
}

func (r *productRepository) FindByID(ctx context.Context, id int) (*domain.Product, error) {
	var model ProductModel
	db := getDB(ctx, r.db)

	if err := db.First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	product := model.ToDomain()
	return &product, nil
}

func (r *productRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Product, error) {
	var model ProductModel
	db := getDB(ctx, r.db)

	if err := db.Where("public_id = ?", publicID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	product := model.ToDomain()
	return &product, nil
}

func (r *productRepository) GetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string) (*domain.Product, error) {
	var model ProductModel
	db := getDB(ctx, r.db)

	// Fetch base product first to get internal ID for collection preloading
	// In a fully optimized version, we'd handle everything in raw SQL, but this is a balanced approach.
	if err := db.Where("public_id = ?", publicID).
		Preload("Categories").
		Preload("Tags").
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Optimized single query for Translation and SEO using LEFT JOIN
	type joinResult struct {
		PTID          int       `gorm:"column:pt_id"`
		PTPublicID    uuid.UUID `gorm:"column:pt_public_id"`
		PTLangCode    string    `gorm:"column:pt_lang_code"`
		PTName        string    `gorm:"column:pt_name"`
		PTDescription string    `gorm:"column:pt_description"`
		PTShortDesc   string    `gorm:"column:pt_short_desc"`
		PSID          int       `gorm:"column:ps_id"`
		PSPublicID    uuid.UUID `gorm:"column:ps_public_id"`
		PSLangCode    string    `gorm:"column:ps_lang_code"`
		PSSlug        string    `gorm:"column:ps_slug"`
		PSMetaTitle   string    `gorm:"column:ps_meta_title"`
		PSMetaDesc    string    `gorm:"column:ps_meta_description"`
		PSCanonical   string    `gorm:"column:ps_canonical_url"`
		PSOGImage     string    `gorm:"column:ps_og_image_url"`
	}

	var res joinResult
	query := `
		SELECT 
			pt.id as pt_id, pt.public_id as pt_public_id, pt.lang_code as pt_lang_code, pt.name as pt_name, pt.description as pt_description, pt.short_desc as pt_short_desc,
			ps.id as ps_id, ps.public_id as ps_public_id, ps.lang_code as ps_lang_code, ps.slug as ps_slug, ps.meta_title as ps_meta_title, ps.meta_description as ps_meta_description, ps.canonical_url as ps_canonical_url, ps.og_image_url as ps_og_image_url
		FROM products p
		LEFT JOIN product_translations pt ON pt.product_id = p.id AND pt.lang_code = ?
		LEFT JOIN product_seo ps ON ps.product_id = p.id AND ps.lang_code = ?
		WHERE p.id = ?
	`

	if err := db.Raw(query, langCode, langCode, model.ID).Scan(&res).Error; err == nil {
		if res.PTID != 0 {
			model.Translation = &ProductTranslationModel{
				BaseModel:   BaseModel{ID: res.PTID, PublicID: res.PTPublicID},
				ProductID:   model.ID,
				LangCode:    res.PTLangCode,
				Name:        res.PTName,
				Description: res.PTDescription,
				ShortDesc:   res.PTShortDesc,
			}
		}
		if res.PSID != 0 {
			model.ProductSEO = &ProductSEOModel{
				BaseModel:       BaseModel{ID: res.PSID, PublicID: res.PSPublicID},
				ProductID:       model.ID,
				LangCode:        res.PSLangCode,
				Slug:            res.PSSlug,
				MetaTitle:       res.PSMetaTitle,
				MetaDescription: res.PSMetaDesc,
				CanonicalURL:    res.PSCanonical,
				OGImageURL:      res.PSOGImage,
			}
		}
	}

	product := model.ToDomain()
	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, p *domain.Product) error {
	model := FromProductDomain(p) // Updated mapper name to FromProductDomain in previous turn
	db := getDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		return err
	}
	// Update domain entity with generated values (like auto-increment ID)
	p.ID = model.ID
	p.CreatedAt = model.CreatedAt
	return nil
}

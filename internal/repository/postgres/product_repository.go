package postgres

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"ss-catalog-service/internal/domain"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		return nil, mapDBError(err)
	}
	product := model.ToDomain()
	return &product, nil
}

func (r *productRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Product, error) {
	var model ProductModel
	db := getDB(ctx, r.db)

	if err := db.Where("public_id = ?", publicID).First(&model).Error; err != nil {
		return nil, mapDBError(err)
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
		return mapDBError(err)
	}
	// Update domain entity with generated values (like auto-increment ID)
	p.ID = model.ID
	p.CreatedAt = model.CreatedAt
	return nil
}

func (r *productRepository) Update(ctx context.Context, p *domain.Product) error {
	model := FromProductDomain(p)
	db := getDB(ctx, r.db)

	// Partial Update: Only update fields that are explicitly provided or changed.
	// Using map with Updates() is the best practice to avoid overwriting all columns.
	updates := map[string]interface{}{
		"name":        model.Name,
		"description": model.Description,
		"status":      model.Status,
		"updated_at":  time.Now(),
	}

	if err := db.Model(&ProductModel{}).Where("id = ?", model.ID).Updates(updates).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (r *productRepository) Search(ctx context.Context, q domain.GetProductSearchQuery) (*domain.ProductSearchResult, error) {
	db := getDB(ctx, r.db)

	// --- Base query scoped to active, non-deleted products ---
	tx := db.Model(&ProductModel{}).
		Where("products.deleted_at IS NULL")

	// --- Full-Text Search via GIN-indexed search_vector ---
	if q.Keyword != nil && *q.Keyword != "" {
		tx = tx.Where("products.search_vector @@ plainto_tsquery('english', ?)", *q.Keyword).
			Order(clause.Expr{
				SQL:  "ts_rank(products.search_vector, plainto_tsquery('english', ?)) DESC",
				Vars: []interface{}{*q.Keyword},
			})
	}

	// --- Status filter (defaults to 'active') ---
	status := domain.ProductStatusActive
	if q.Status != nil {
		status = *q.Status
	}
	tx = tx.Where("products.status = ?", status)

	// --- Brand filter ---
	if q.BrandID != nil {
		tx = tx.Where("products.brand_id = ?", *q.BrandID)
	}

	// --- Category filter via JOIN ---
	if q.CategorySlug != nil && *q.CategorySlug != "" {
		tx = tx.Joins(`
			INNER JOIN product_categories pc ON pc.product_id = products.id
			INNER JOIN categories c ON c.id = pc.category_id AND c.deleted_at IS NULL
		`).Where("c.slug = ?", *q.CategorySlug)
	}

	// --- Price range filter via subquery on variants ---
	if q.MinPrice != nil || q.MaxPrice != nil {
		tx = tx.Joins(`
			INNER JOIN (
				SELECT product_id, MIN(price) AS min_price
				FROM product_variants
				WHERE deleted_at IS NULL
				GROUP BY product_id
			) pv ON pv.product_id = products.id
		`)
		if q.MinPrice != nil {
			tx = tx.Where("pv.min_price >= ?", *q.MinPrice)
		}
		if q.MaxPrice != nil {
			tx = tx.Where("pv.min_price <= ?", *q.MaxPrice)
		}
	}

	// --- Cursor-based pagination ---
	if q.Cursor != nil && *q.Cursor != "" {
		cursorID, cursorCreatedAt, err := decodeCursor(*q.Cursor)
		if err != nil {
			return nil, domain.ErrInvalidCursor
		}
		tx = tx.Where(
			"(products.created_at, products.id) < (?, ?)",
			cursorCreatedAt, cursorID,
		)
	}

	// --- Keyset sort + limit ---
	tx = tx.Order("products.created_at DESC, products.id DESC").Limit(q.Limit + 1)

	// --- Execute ---
	var models []ProductModel
	if err := tx.Find(&models).Error; err != nil {
		return nil, mapDBError(err)
	}

	// --- Build result + next cursor ---
	result := &domain.ProductSearchResult{}

	// --- Implement TotalHint using approximate count ---
	// Using reltuples for performance (O(1) vs O(N) for COUNT(*))
	var hint int64
	r.db.Raw("SELECT reltuples::bigint AS estimate FROM pg_class WHERE relname = 'products'").Scan(&hint)
	result.TotalHint = hint

	if len(models) > q.Limit {
		models = models[:q.Limit]
		last := models[len(models)-1]
		cursor := encodeCursor(last.ID, last.CreatedAt)
		result.NextCursor = &cursor
	}

	result.Items = make([]domain.Product, len(models))
	for i, m := range models {
		result.Items[i] = m.ToDomain()
	}

	return result, nil
}

func encodeCursor(id int, createdAt time.Time) string {
	raw := fmt.Sprintf("%d:%d", id, createdAt.UnixNano())
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func decodeCursor(cursor string) (int, time.Time, error) {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, time.Time{}, err
	}
	parts := strings.SplitN(string(b), ":", 2)
	if len(parts) != 2 {
		return 0, time.Time{}, errors.New("malformed cursor")
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, time.Time{}, err
	}
	nano, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, time.Time{}, err
	}
	return id, time.Unix(0, nano), nil
}

// mapDBError translates database-specific errors into domain errors.
func mapDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrProductNotFound
	}

	// Check for Postgres Unique Constraint Violation (Code 23505)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrDuplicateProduct
	}

	// General GORM or connection errors
	if strings.Contains(err.Error(), "duplicate key value") {
		return domain.ErrDuplicateProduct
	}

	return fmt.Errorf("%w: %v", domain.ErrInternalDatabase, err)
}

package domain

// ProductCategory maps a product to its categories (M2M).
type ProductCategory struct {
	BaseEntity
	ProductID  int
	CategoryID int
	IsPrimary  bool
}

// Tag represents flat keyword tags for products.
type Tag struct {
	BaseEntity
	Name string
	Slug string
}

// ProductTag maps a product to its tags (M2M).
type ProductTag struct {
	BaseEntity
	ProductID int
	TagID     int
}

// SEOBase contains common SEO fields.
type SEOBase struct {
	LangCode        string
	Slug            string
	MetaTitle       string
	MetaDescription string
}

// ProductSEO represents SEO metadata for products.
type ProductSEO struct {
	BaseEntity
	ProductID int
	SEOBase
	CanonicalURL string
	OGImageURL   string
}

// CategorySEO represents SEO metadata for categories.
type CategorySEO struct {
	BaseEntity
	CategoryID int
	SEOBase
}

// ProductTranslation represents localized product data.
type ProductTranslation struct {
	BaseEntity
	ProductID int
	LangCode  string
	Name      string
	Description string
	ShortDesc   string
}

// CategoryTranslation represents localized category names.
type CategoryTranslation struct {
	BaseEntity
	CategoryID int
	LangCode   string
	Name       string
	Description string
}

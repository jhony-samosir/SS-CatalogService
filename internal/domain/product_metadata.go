package domain

// ProductCategory maps a product to its categories (M2M).
type ProductCategory struct {
	BaseEntity
	ProductID  int  `json:"product_id"`
	CategoryID int  `json:"category_id"`
	IsPrimary  bool `json:"is_primary"`
}

// Tag represents flat keyword tags for products.
type Tag struct {
	BaseEntity
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ProductTag maps a product to its tags (M2M).
type ProductTag struct {
	BaseEntity
	ProductID int `json:"product_id"`
	TagID     int `json:"tag_id"`
}

// SEOBase contains common SEO fields.
type SEOBase struct {
	LangCode        string `json:"lang_code"`
	Slug            string `json:"slug"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}

// ProductSEO represents SEO metadata for products.
type ProductSEO struct {
	BaseEntity
	ProductID int `json:"product_id"`
	SEOBase
	CanonicalURL string `json:"canonical_url"`
	OGImageURL   string `json:"og_image_url"`
}

// CategorySEO represents SEO metadata for categories.
type CategorySEO struct {
	BaseEntity
	CategoryID int `json:"category_id"`
	SEOBase
}

// ProductTranslation represents localized product data.
type ProductTranslation struct {
	BaseEntity
	ProductID int    `json:"product_id"`
	LangCode  string `json:"lang_code"`
	Name      string `json:"name"`
	Description string `json:"description"`
	ShortDesc   string `json:"short_desc"`
}

// CategoryTranslation represents localized category names.
type CategoryTranslation struct {
	BaseEntity
	CategoryID int    `json:"category_id"`
	LangCode   string `json:"lang_code"`
	Name       string `json:"name"`
	Description string `json:"description"`
}

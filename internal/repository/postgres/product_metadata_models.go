package postgres

// ProductCategoryModel represents the database schema for product_categories.
type ProductCategoryModel struct {
	BaseModel
	ProductID  int  `gorm:"not null;index"`
	CategoryID int  `gorm:"not null;index"`
	IsPrimary  bool `gorm:"not null;default:false"`
}

func (ProductCategoryModel) TableName() string { return "product_categories" }

// TagModel represents the database schema for tags.
type TagModel struct {
	BaseModel
	Name string `gorm:"type:varchar(100);not null"`
	Slug string `gorm:"type:varchar(100);not null;uniqueIndex"`
}

func (TagModel) TableName() string { return "tags" }

// ProductTagModel represents the database schema for product_tags.
type ProductTagModel struct {
	BaseModel
	ProductID int `gorm:"not null;index"`
	TagID     int `gorm:"not null;index"`
}

func (ProductTagModel) TableName() string { return "product_tags" }

// ProductSEOModel represents the database schema for product_seo.
type ProductSEOModel struct {
	BaseModel
	ProductID       int    `gorm:"not null;index"`
	LangCode        string `gorm:"type:char(5);not null;default:'id-ID'"`
	Slug            string `gorm:"type:varchar(500);not null;uniqueIndex"`
	MetaTitle       string `gorm:"type:varchar(255)"`
	MetaDescription string `gorm:"type:varchar(500)"`
	CanonicalURL    string `gorm:"type:text"`
	OGImageURL      string `gorm:"type:text"`
}

func (ProductSEOModel) TableName() string { return "product_seo" }

// CategorySEOModel represents the database schema for category_seo.
type CategorySEOModel struct {
	BaseModel
	CategoryID      int    `gorm:"not null;index"`
	LangCode        string `gorm:"type:char(5);not null;default:'id-ID'"`
	Slug            string `gorm:"type:varchar(500);not null;uniqueIndex"`
	MetaTitle       string `gorm:"type:varchar(255)"`
	MetaDescription string `gorm:"type:varchar(500)"`
}

func (CategorySEOModel) TableName() string { return "category_seo" }

// ProductTranslationModel represents the database schema for product_translations.
type ProductTranslationModel struct {
	BaseModel
	ProductID   int    `gorm:"not null;index"`
	LangCode    string `gorm:"type:char(5);not null"`
	Name        string `gorm:"type:varchar(500);not null"`
	Description string `gorm:"type:text"`
	ShortDesc   string `gorm:"type:varchar(1000)"`
}

func (ProductTranslationModel) TableName() string { return "product_translations" }

// CategoryTranslationModel represents the database schema for category_translations.
type CategoryTranslationModel struct {
	BaseModel
	CategoryID  int    `gorm:"not null;index"`
	LangCode    string `gorm:"type:char(5);not null"`
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
}

func (CategoryTranslationModel) TableName() string { return "category_translations" }

package domain

// ProductImage represents product and variant images.
type ProductImage struct {
	BaseEntity
	ProductID  *int
	VariantID  *int
	URL        string
	AltText    string
	SortOrder  int
	IsPrimary  bool
	WidthPX    int
	HeightPX   int
}

// ProductVideo represents product demo/promo videos.
type ProductVideo struct {
	BaseEntity
	ProductID   int
	URL         string
	Thumbnail   string
	Title       string
	DurationSec int
	SortOrder   int
}

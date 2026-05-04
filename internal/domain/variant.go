package domain

// ProductVariant represents a SKU-level sellable product unit.
type ProductVariant struct {
	BaseEntity
	ProductID  int
	SKU        string
	Barcode    string
	Name       string
	IsDefault  bool
	IsActive   bool
	WeightGram *int
	SortOrder  int
}

package domain

import "time"

// PriceType defines the tier of pricing.
type PriceType string

const (
	PriceTypeRetail    PriceType = "retail"
	PriceTypeDiscount  PriceType = "discount"
	PriceTypeWholesale PriceType = "wholesale"
	PriceTypeMember    PriceType = "member"
)

// ProductPrice represents multi-type pricing per variant.
type ProductPrice struct {
	BaseEntity
	VariantID    int
	PriceType    PriceType
	CurrencyCode string
	Amount       float64
	MinQuantity  int
	ValidFrom    *time.Time
	ValidUntil   *time.Time
	IsActive     bool
}

// Warehouse represents physical or virtual warehouse locations.
type Warehouse struct {
	BaseEntity
	SellerID    *int
	Name        string
	Code        string
	City        string
	Province    string
	CountryCode string
	PostalCode  string
	Address     string
	IsActive    bool
}

// ProductInventory represents real-time stock levels.
type ProductInventory struct {
	BaseEntity
	VariantID        int
	WarehouseID      int
	QuantityOnHand   int
	QuantityReserved int
	LowStockAlert    int
}

// MovementType defines the type of stock change.
type MovementType string

const (
	MovementTypeIn         MovementType = "in"
	MovementTypeOut        MovementType = "out"
	MovementTypeReserve    MovementType = "reserve"
	MovementTypeRelease    MovementType = "release"
	MovementTypeAdjustment MovementType = "adjustment"
)

// InventoryMovement represents an audit trail for stock changes.
type InventoryMovement struct {
	BaseEntity
	InventoryID   int
	MovementType  MovementType
	Quantity      int
	ReferenceID   string
	ReferenceType string
	Note          string
}

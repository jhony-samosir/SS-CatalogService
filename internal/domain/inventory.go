package domain

import (
	"context"
	"errors"
	"time"
)

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

// InventoryReferenceType defines standardized sources of stock changes.
type InventoryReferenceType string

const (
	ReferenceTypeOrder      InventoryReferenceType = "order"
	ReferenceTypeRestock    InventoryReferenceType = "restock"
	ReferenceTypeAdjustment InventoryReferenceType = "adjustment"
	ReferenceTypeReturn     InventoryReferenceType = "return"
)

// InventoryMovement represents an audit trail for stock changes.
type InventoryMovement struct {
	BaseEntity
	InventoryID   int
	MovementType  MovementType
	Quantity      int
	ReferenceID   string
	ReferenceType InventoryReferenceType
	Note          string
}

// Sentinel Errors
var (
	ErrInsufficientStock = errors.New("insufficient stock available")
	ErrInventoryNotFound = errors.New("inventory record not found")
)

// --- DTOs ---

type UpdateStockPayload struct {
	VariantID     int
	WarehouseID   int
	Quantity      int // Positive for IN, Negative for OUT
	ReferenceType InventoryReferenceType
	ReferenceID   string
	Note          string
}

// --- Interfaces ---

type InventoryRepository interface {
	GetInventoryForUpdate(ctx context.Context, variantID int, warehouseID int) (*ProductInventory, error)
	CreateInventory(ctx context.Context, inv *ProductInventory) error
	UpdateInventory(ctx context.Context, inv *ProductInventory) error
	CreateMovement(ctx context.Context, movement *InventoryMovement) error
}

type InventoryCommandUsecase interface {
	UpdateInventoryStock(ctx context.Context, payload UpdateStockPayload) error
}

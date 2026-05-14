package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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
	SellerID    *int   `json:"seller_id,omitempty"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	City        string `json:"city"`
	Province    string `json:"province,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Address     string `json:"address,omitempty"`
	IsActive    bool   `json:"is_active"`
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

type WarehouseRepository interface {
	FindAll(ctx context.Context, p Pagination) ([]Warehouse, error)
	Count(ctx context.Context) (int64, error)
	FindByPublicID(ctx context.Context, publicID uuid.UUID) (*Warehouse, error)
	Create(ctx context.Context, wh *Warehouse) error
	Update(ctx context.Context, wh *Warehouse) error
	Delete(ctx context.Context, publicID uuid.UUID) error
	CountInventory(ctx context.Context, warehouseID int) (int64, error)
}

type InventoryCommandUsecase interface {
	UpdateInventoryStock(ctx context.Context, payload UpdateStockPayload) error
}

// WarehouseUsecase defines the business logic for warehouses.
type WarehouseUsecase interface {
	GetWarehouses(ctx context.Context, p Pagination) ([]Warehouse, int64, error)
	CreateWarehouse(ctx context.Context, wh *Warehouse) error
	UpdateWarehouse(ctx context.Context, wh *Warehouse) error
	DeleteWarehouse(ctx context.Context, publicID uuid.UUID) error
}

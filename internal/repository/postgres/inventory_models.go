package postgres

import (
	"time"
)

// ProductPriceModel represents the database schema for product_prices.
type ProductPriceModel struct {
	BaseModel
	VariantID    int       `gorm:"not null;index"`
	PriceType    string    `gorm:"type:varchar(50);not null;default:'retail'"`
	CurrencyCode string    `gorm:"type:char(3);not null;default:'IDR'"`
	Amount       float64   `gorm:"type:numeric(18,2);not null"`
	MinQuantity  int       `gorm:"not null;default:1"`
	ValidFrom    *time.Time `gorm:"null"`
	ValidUntil   *time.Time `gorm:"null"`
	IsActive     bool      `gorm:"not null;default:true"`
}

func (ProductPriceModel) TableName() string { return "product_prices" }

// WarehouseModel represents the database schema for warehouses.
type WarehouseModel struct {
	BaseModel
	SellerID    *int   `gorm:"index"`
	Name        string `gorm:"type:varchar(255);not null"`
	Code        string `gorm:"type:varchar(100);not null;uniqueIndex"`
	City        string `gorm:"type:varchar(255)"`
	Province    string `gorm:"type:varchar(255)"`
	CountryCode string `gorm:"type:char(2);not null;default:'ID'"`
	PostalCode  string `gorm:"type:varchar(20)"`
	Address     string `gorm:"type:text"`
	IsActive    bool   `gorm:"not null;default:true"`
}

func (WarehouseModel) TableName() string { return "warehouses" }

// ProductInventoryModel represents the database schema for product_inventory.
type ProductInventoryModel struct {
	BaseModel
	VariantID        int `gorm:"not null;index"`
	WarehouseID      int `gorm:"not null;index"`
	QuantityOnHand   int `gorm:"not null;default:0"`
	QuantityReserved int `gorm:"not null;default:0"`
	LowStockAlert    int `gorm:"not null;default:5"`
}

func (ProductInventoryModel) TableName() string { return "product_inventory" }

// InventoryMovementModel represents the database schema for inventory_movements.
type InventoryMovementModel struct {
	BaseModel
	InventoryID   int    `gorm:"not null;index"`
	MovementType  string `gorm:"type:varchar(50);not null"`
	Quantity      int    `gorm:"not null"`
	ReferenceID   string `gorm:"type:varchar(255);index"`
	ReferenceType string `gorm:"type:varchar(100);index"`
	Note          string `gorm:"type:text"`
}

func (InventoryMovementModel) TableName() string { return "inventory_movements" }

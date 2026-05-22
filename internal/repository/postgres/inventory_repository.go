package postgres

import (
	"context"
	"errors"
	"log/slog"
	"ss-catalog-service/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) domain.InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) FindAll(ctx context.Context, p domain.Pagination, warehouseID string, variantID string) ([]domain.ProductInventory, int64, error) {
	var results []domain.ProductInventory
	var total int64
	db := getDB(ctx, r.db)

	query := db.Table("product_inventory").
		Select(`
			product_inventory.id, 
			product_inventory.public_id, 
			product_inventory.created_at, 
			product_inventory.updated_at,
			product_inventory.variant_id, 
			product_inventory.warehouse_id, 
			product_inventory.quantity_on_hand, 
			product_inventory.quantity_reserved, 
			product_inventory.low_stock_alert,
			products.name as product_name, 
			product_variants.name as variant_name, 
			product_variants.sku, 
			warehouses.name as warehouse_name
		`).
		Joins("JOIN product_variants ON product_variants.id = product_inventory.variant_id").
		Joins("JOIN products ON products.id = product_variants.product_id").
		Joins("JOIN warehouses ON warehouses.id = product_inventory.warehouse_id").
		Where("product_inventory.deleted_at IS NULL")

	if warehouseID != "" {
		query = query.Where("product_inventory.warehouse_id = ?", warehouseID)
	}
	if variantID != "" {
		query = query.Where("product_inventory.variant_id = ?", variantID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if p.Limit > 0 {
		query = query.Limit(p.Limit).Offset(p.Offset)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *inventoryRepository) GetInventoryForUpdate(ctx context.Context, variantID int, warehouseID int) (*domain.ProductInventory, error) {
	var model ProductInventoryModel
	db := getDB(ctx, r.db)

	// Best Practice: Use FOR UPDATE to lock the row for the duration of the transaction
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("variant_id = ? AND warehouse_id = ?", variantID, warehouseID).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	inv := model.ToDomain()
	return &inv, nil
}

func (r *inventoryRepository) CreateInventory(ctx context.Context, inv *domain.ProductInventory) error {
	model := FromInventoryDomain(inv)
	db := getDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "create_inventory", "error", err)
		return err
	}
	inv.ID = model.ID
	inv.CreatedAt = model.CreatedAt
	inv.PublicID = model.PublicID
	slog.InfoContext(ctx, "db_audit", "operation", "create_inventory", "inventory_id", inv.ID, "variant_id", inv.VariantID, "warehouse_id", inv.WarehouseID)
	return nil
}

func (r *inventoryRepository) UpdateInventory(ctx context.Context, inv *domain.ProductInventory) error {
	model := FromInventoryDomain(inv)
	db := getDB(ctx, r.db)

	// Best Practice: Partial update to prevent overwriting audit fields like CreatedAt
	err := db.Model(model).Select("QuantityOnHand", "QuantityReserved", "UpdatedAt", "UpdatedBy").Updates(model).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "update_inventory", "inventory_id", inv.ID, "error", err)
		return err
	}
	slog.InfoContext(ctx, "db_audit", "operation", "update_inventory", "inventory_id", inv.ID, "qty_on_hand", inv.QuantityOnHand, "qty_reserved", inv.QuantityReserved)
	return nil
}

func (r *inventoryRepository) CreateMovement(ctx context.Context, movement *domain.InventoryMovement) error {
	model := FromMovementDomain(movement)
	db := getDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "create_movement", "error", err)
		return err
	}
	movement.ID = model.ID
	slog.InfoContext(ctx, "db_audit", "operation", "create_movement", "movement_id", movement.ID, "inventory_id", movement.InventoryID, "quantity", movement.Quantity, "movement_type", movement.MovementType)
	return nil
}

package postgres

import (
	"context"
	"errors"
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
		return err
	}
	inv.ID = model.ID
	inv.CreatedAt = model.CreatedAt
	inv.PublicID = model.PublicID
	return nil
}

func (r *inventoryRepository) UpdateInventory(ctx context.Context, inv *domain.ProductInventory) error {
	model := FromInventoryDomain(inv)
	db := getDB(ctx, r.db)

	// Best Practice: Partial update to prevent overwriting audit fields like CreatedAt
	return db.Model(model).Select("QuantityOnHand", "QuantityReserved", "UpdatedAt").Updates(model).Error
}

func (r *inventoryRepository) CreateMovement(ctx context.Context, movement *domain.InventoryMovement) error {
	model := FromMovementDomain(movement)
	db := getDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		return err
	}
	movement.ID = model.ID
	return nil
}

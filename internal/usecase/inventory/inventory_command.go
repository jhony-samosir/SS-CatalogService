package inventory

import (
	"context"
	"fmt"
	"log/slog"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type inventoryCommandUsecase struct {
	invRepo   domain.InventoryRepository
	txManager domain.TransactionManager
	logger    *slog.Logger
}

func NewInventoryCommandUsecase(repo domain.InventoryRepository, txManager domain.TransactionManager, logger *slog.Logger) domain.InventoryCommandUsecase {
	return &inventoryCommandUsecase{
		invRepo:   repo,
		txManager: txManager,
		logger:    logger,
	}
}

func (u *inventoryCommandUsecase) UpdateInventoryStock(ctx context.Context, payload domain.UpdateStockPayload) error {
	// Best Practice 2: Zero Quantity Guard Clause (Save DB resources)
	if payload.Quantity == 0 {
		return nil
	}

	u.logger.InfoContext(ctx, "updating inventory stock start", 
		"variant_id", payload.VariantID, 
		"warehouse_id", payload.WarehouseID, 
		"quantity", payload.Quantity,
		"ref_type", payload.ReferenceType,
		"ref_id", payload.ReferenceID,
	)

	var finalQty int
	// Execute within a single database transaction
	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Lock and Read Current Inventory
		inv, err := u.invRepo.GetInventoryForUpdate(txCtx, payload.VariantID, payload.WarehouseID)
		if err != nil {
			u.logger.ErrorContext(txCtx, "failed to get inventory for update", "error", err, "variant_id", payload.VariantID, "warehouse_id", payload.WarehouseID)
			return fmt.Errorf("UpdateInventoryStock.GetInventory: %w", err)
		}

		// Best Practice 3: Handle Missing Inventory (Auto-Create for IN, Error for OUT)
		if inv == nil {
			if payload.Quantity < 0 {
				u.logger.WarnContext(txCtx, "insufficient stock: inventory record not found", "variant_id", payload.VariantID, "warehouse_id", payload.WarehouseID)
				return domain.ErrInventoryNotFound
			}

			// Auto-Create for first time restock/adjustment
			inv = &domain.ProductInventory{
				BaseEntity:     domain.BaseEntity{PublicID: uuid.New()},
				VariantID:      payload.VariantID,
				WarehouseID:    payload.WarehouseID,
				QuantityOnHand: 0,
			}
			u.logger.InfoContext(txCtx, "auto-creating product inventory record", "variant_id", payload.VariantID, "warehouse_id", payload.WarehouseID)
			if err := u.invRepo.CreateInventory(txCtx, inv); err != nil {
				u.logger.ErrorContext(txCtx, "failed to auto-create inventory record", "error", err, "variant_id", payload.VariantID, "warehouse_id", payload.WarehouseID)
				return fmt.Errorf("UpdateInventoryStock.AutoCreateInventory: %w", err)
			}
		}

		// 2. Business Logic / Validation
		if payload.Quantity < 0 && inv.QuantityOnHand < (-payload.Quantity) {
			u.logger.WarnContext(txCtx, "insufficient stock validation failed", 
				"variant_id", payload.VariantID, 
				"warehouse_id", payload.WarehouseID, 
				"on_hand", inv.QuantityOnHand, 
				"requested", -payload.Quantity,
			)
			return domain.ErrInsufficientStock
		}

		// Update stock locally
		inv.QuantityOnHand += payload.Quantity
		finalQty = inv.QuantityOnHand

		// 3. Create Inventory Movement (Ledger)
		movementType := domain.MovementTypeIn
		if payload.Quantity < 0 {
			movementType = domain.MovementTypeOut
		}

		movement := &domain.InventoryMovement{
			BaseEntity:    domain.BaseEntity{PublicID: uuid.New()},
			InventoryID:   inv.ID,
			MovementType:  movementType,
			Quantity:      payload.Quantity,
			ReferenceID:   payload.ReferenceID,
			ReferenceType: payload.ReferenceType,
			Note:          payload.Note,
		}

		if err := u.invRepo.CreateMovement(txCtx, movement); err != nil {
			u.logger.ErrorContext(txCtx, "failed to create inventory movement log", "error", err, "inventory_id", inv.ID)
			return fmt.Errorf("UpdateInventoryStock.CreateMovement: %w", err)
		}

		// 4. Update Product Inventory
		if err := u.invRepo.UpdateInventory(txCtx, inv); err != nil {
			u.logger.ErrorContext(txCtx, "failed to update product inventory", "error", err, "inventory_id", inv.ID)
			return fmt.Errorf("UpdateInventoryStock.UpdateInventory: %w", err)
		}

		return nil
	})

	if err != nil {
		u.logger.ErrorContext(ctx, "inventory stock update transaction failed", "error", err, "variant_id", payload.VariantID, "warehouse_id", payload.WarehouseID)
		return fmt.Errorf("UpdateInventoryStock transaction failed: %w", err)
	}

	u.logger.InfoContext(ctx, "inventory stock updated successfully", 
		"variant_id", payload.VariantID, 
		"warehouse_id", payload.WarehouseID, 
		"final_quantity", finalQty,
	)
	return nil
}

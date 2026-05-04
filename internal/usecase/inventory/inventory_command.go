package inventory

import (
	"context"
	"fmt"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type inventoryCommandUsecase struct {
	invRepo   domain.InventoryRepository
	txManager domain.TransactionManager
}

func NewInventoryCommandUsecase(repo domain.InventoryRepository, txManager domain.TransactionManager) domain.InventoryCommandUsecase {
	return &inventoryCommandUsecase{
		invRepo:   repo,
		txManager: txManager,
	}
}

func (u *inventoryCommandUsecase) UpdateInventoryStock(ctx context.Context, payload domain.UpdateStockPayload) error {
	// Best Practice 2: Zero Quantity Guard Clause (Save DB resources)
	if payload.Quantity == 0 {
		return nil
	}

	// Execute within a single database transaction
	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Lock and Read Current Inventory
		inv, err := u.invRepo.GetInventoryForUpdate(txCtx, payload.VariantID, payload.WarehouseID)
		if err != nil {
			return fmt.Errorf("UpdateInventoryStock.GetInventory: %w", err)
		}

		// Best Practice 3: Handle Missing Inventory (Auto-Create for IN, Error for OUT)
		if inv == nil {
			if payload.Quantity < 0 {
				return domain.ErrInventoryNotFound
			}

			// Auto-Create for first time restock/adjustment
			inv = &domain.ProductInventory{
				BaseEntity:     domain.BaseEntity{PublicID: uuid.New()},
				VariantID:      payload.VariantID,
				WarehouseID:    payload.WarehouseID,
				QuantityOnHand: 0,
			}
			if err := u.invRepo.CreateInventory(txCtx, inv); err != nil {
				return fmt.Errorf("UpdateInventoryStock.AutoCreateInventory: %w", err)
			}
		}

		// 2. Business Logic / Validation
		if payload.Quantity < 0 && inv.QuantityOnHand < (-payload.Quantity) {
			return domain.ErrInsufficientStock
		}

		// Update stock locally
		inv.QuantityOnHand += payload.Quantity

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
			return fmt.Errorf("UpdateInventoryStock.CreateMovement: %w", err)
		}

		// 4. Update Product Inventory
		if err := u.invRepo.UpdateInventory(txCtx, inv); err != nil {
			return fmt.Errorf("UpdateInventoryStock.UpdateInventory: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("UpdateInventoryStock transaction failed: %w", err)
	}

	return nil
}

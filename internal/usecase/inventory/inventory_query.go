package inventory

import (
	"context"
	"ss-catalog-service/internal/domain"
)

type inventoryQueryUsecase struct {
	repo domain.InventoryRepository
}

func NewInventoryQueryUsecase(repo domain.InventoryRepository) domain.InventoryQueryUsecase {
	return &inventoryQueryUsecase{repo: repo}
}

func (u *inventoryQueryUsecase) GetInventory(ctx context.Context, p domain.Pagination, warehouseID string, variantID string) ([]domain.ProductInventory, int64, error) {
	return u.repo.FindAll(ctx, p, warehouseID, variantID)
}

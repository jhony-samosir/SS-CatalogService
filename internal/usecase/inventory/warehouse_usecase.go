package inventory

import (
	"context"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type warehouseUsecase struct {
	repo domain.WarehouseRepository
}

func NewWarehouseUsecase(repo domain.WarehouseRepository) domain.WarehouseUsecase {
	return &warehouseUsecase{repo: repo}
}

func (u *warehouseUsecase) GetWarehouses(ctx context.Context, p domain.Pagination) ([]domain.Warehouse, error) {
	return u.repo.FindAll(ctx, p)
}

func (u *warehouseUsecase) CreateWarehouse(ctx context.Context, wh *domain.Warehouse) error {
	if wh.PublicID == uuid.Nil {
		wh.PublicID = uuid.New()
	}
	return u.repo.Create(ctx, wh)
}

func (u *warehouseUsecase) UpdateWarehouse(ctx context.Context, wh *domain.Warehouse) error {
	return u.repo.Update(ctx, wh)
}

func (u *warehouseUsecase) DeleteWarehouse(ctx context.Context, publicID uuid.UUID) error {
	return u.repo.Delete(ctx, publicID)
}

package product

import (
	"context"
	"fmt"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type productQueryUsecase struct {
	repo domain.ProductRepository
}

// NewProductQueryUsecase creates a new instance of product query business logic.
func NewProductQueryUsecase(repo domain.ProductRepository) domain.ProductQueryUsecase {
	return &productQueryUsecase{repo: repo}
}

func (u *productQueryUsecase) GetAllProducts(ctx context.Context, p domain.Pagination) ([]domain.Product, error) {
	products, err := u.repo.FindAll(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("productQueryUsecase.GetAllProducts: %w", err)
	}
	return products, nil
}

func (u *productQueryUsecase) GetProductByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Product, error) {
	product, err := u.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return nil, fmt.Errorf("productQueryUsecase.GetProductByPublicID: %w", err)
	}
	if product == nil {
		return nil, domain.ErrProductNotFound
	}
	return product, nil
}

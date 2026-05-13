package category

import (
	"context"
	"fmt"
	"ss-catalog-service/internal/domain"
)

type categoryUsecase struct {
	repo domain.CategoryRepository
}

// NewCategoryUsecase creates a new instance of category business logic.
func NewCategoryUsecase(repo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		repo: repo,
	}
}

func (u *categoryUsecase) GetCategories(ctx context.Context, p domain.Pagination) ([]domain.Category, error) {
	categories, err := u.repo.FindAll(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("categoryUsecase.GetCategories: %w", err)
	}
	return categories, nil
}

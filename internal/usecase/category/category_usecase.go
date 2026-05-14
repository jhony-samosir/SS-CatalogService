package category

import (
	"context"
	"fmt"
	"ss-catalog-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type categoryUsecase struct {
	repo  domain.CategoryRepository
	cache domain.MasterDataCacheRepository
}

func NewCategoryUsecase(repo domain.CategoryRepository, cache domain.MasterDataCacheRepository) domain.CategoryUsecase {
	return &categoryUsecase{repo: repo, cache: cache}
}

func (u *categoryUsecase) GetCategories(ctx context.Context, p domain.Pagination) ([]domain.Category, error) {
	cacheKey := fmt.Sprintf("categories:all:%d:%d", p.Limit, p.Offset)
	var categories []domain.Category

	if err := u.cache.Get(ctx, cacheKey, &categories); err == nil {
		return categories, nil
	}

	categories, err := u.repo.FindAll(ctx, p)
	if err != nil {
		return nil, err
	}

	_ = u.cache.Set(ctx, cacheKey, categories, 1*time.Hour)
	return categories, nil
}

func (u *categoryUsecase) GetCategoryByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Category, error) {
	cacheKey := fmt.Sprintf("category:%s", publicID.String())
	var category domain.Category

	if err := u.cache.Get(ctx, cacheKey, &category); err == nil {
		return &category, nil
	}

	res, err := u.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	_ = u.cache.Set(ctx, cacheKey, res, 1*time.Hour)
	return res, nil
}

func (u *categoryUsecase) CreateCategory(ctx context.Context, category *domain.Category) error {
	if category.PublicID == uuid.Nil {
		category.PublicID = uuid.New()
	}
	if err := u.repo.Create(ctx, category); err != nil {
		return err
	}

	_ = u.cache.InvalidateAll(ctx)
	return nil
}

func (u *categoryUsecase) UpdateCategory(ctx context.Context, category *domain.Category) error {
	if err := u.repo.Update(ctx, category); err != nil {
		return err
	}

	_ = u.cache.InvalidateAll(ctx)
	_ = u.cache.Delete(ctx, fmt.Sprintf("category:%s", category.PublicID.String()))
	return nil
}

func (u *categoryUsecase) DeleteCategory(ctx context.Context, publicID uuid.UUID) error {
	category, err := u.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return err
	}
	if category == nil {
		return domain.ErrNotFound
	}

	// Check for children
	childCount, err := u.repo.CountChildren(ctx, category.ID)
	if err != nil {
		return err
	}
	if childCount > 0 {
		return domain.ErrEntityInUse
	}

	// Check for products
	prodCount, err := u.repo.CountProducts(ctx, category.ID)
	if err != nil {
		return err
	}
	if prodCount > 0 {
		return domain.ErrEntityInUse
	}

	if err := u.repo.Delete(ctx, publicID); err != nil {
		return err
	}

	_ = u.cache.InvalidateAll(ctx)
	_ = u.cache.Delete(ctx, fmt.Sprintf("category:%s", publicID.String()))
	return nil
}

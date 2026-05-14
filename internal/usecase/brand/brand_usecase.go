package brand

import (
	"context"
	"fmt"
	"ss-catalog-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type brandUsecase struct {
	repo  domain.BrandRepository
	cache domain.MasterDataCacheRepository
}

func NewBrandUsecase(repo domain.BrandRepository, cache domain.MasterDataCacheRepository) domain.BrandUsecase {
	return &brandUsecase{repo: repo, cache: cache}
}

func (u *brandUsecase) GetBrands(ctx context.Context, p domain.Pagination) ([]domain.Brand, int64, error) {
	cacheKey := fmt.Sprintf("brands:all:%d:%d", p.Limit, p.Offset)
	var brands []domain.Brand
	var total int64

	// For simplicity in this session, we'll cache the result but we need the total too.
	// In a real scenario, we might cache the total count separately or wrap it.
	if err := u.cache.Get(ctx, cacheKey, &brands); err == nil {
		// If we hit cache, we still need the total count.
		// For now, let's just fetch total from repo (or cache it too).
		total, _ = u.repo.Count(ctx)
		return brands, total, nil
	}

	brands, err := u.repo.FindAll(ctx, p)
	if err != nil {
		return nil, 0, err
	}

	total, err = u.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	_ = u.cache.Set(ctx, cacheKey, brands, 1*time.Hour)
	return brands, total, nil
}

func (u *brandUsecase) GetBrandByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Brand, error) {
	cacheKey := fmt.Sprintf("brand:%s", publicID.String())
	var brand domain.Brand

	if err := u.cache.Get(ctx, cacheKey, &brand); err == nil {
		return &brand, nil
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

func (u *brandUsecase) CreateBrand(ctx context.Context, brand *domain.Brand) error {
	if brand.PublicID == uuid.Nil {
		brand.PublicID = uuid.New()
	}
	if err := u.repo.Create(ctx, brand); err != nil {
		return err
	}

	_ = u.cache.InvalidateAll(ctx)
	return nil
}

func (u *brandUsecase) UpdateBrand(ctx context.Context, brand *domain.Brand) error {
	if err := u.repo.Update(ctx, brand); err != nil {
		return err
	}

	_ = u.cache.InvalidateAll(ctx)
	_ = u.cache.Delete(ctx, fmt.Sprintf("brand:%s", brand.PublicID.String()))
	return nil
}

func (u *brandUsecase) DeleteBrand(ctx context.Context, publicID uuid.UUID) error {
	brand, err := u.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return err
	}
	if brand == nil {
		return domain.ErrNotFound
	}

	count, err := u.repo.CountProducts(ctx, brand.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		return domain.ErrEntityInUse
	}

	if err := u.repo.Delete(ctx, publicID); err != nil {
		return err
	}

	_ = u.cache.InvalidateAll(ctx)
	_ = u.cache.Delete(ctx, fmt.Sprintf("brand:%s", publicID.String()))
	return nil
}

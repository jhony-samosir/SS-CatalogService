package bundle

import (
	"context"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type bundleUsecase struct {
	repo domain.BundleRepository
}

func NewBundleUsecase(repo domain.BundleRepository) domain.BundleUsecase {
	return &bundleUsecase{repo: repo}
}

func (u *bundleUsecase) CreateBundle(ctx context.Context, bundle *domain.ProductBundle) error {
	return u.repo.Create(ctx, bundle)
}

func (u *bundleUsecase) GetBundles(ctx context.Context, p domain.Pagination) ([]domain.ProductBundle, int64, error) {
	return u.repo.FindAll(ctx, p)
}

func (u *bundleUsecase) GetBundleByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.ProductBundle, error) {
	return u.repo.FindByPublicID(ctx, publicID)
}

func (u *bundleUsecase) UpdateBundle(ctx context.Context, bundle *domain.ProductBundle) error {
	return u.repo.Update(ctx, bundle)
}

func (u *bundleUsecase) DeleteBundle(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}

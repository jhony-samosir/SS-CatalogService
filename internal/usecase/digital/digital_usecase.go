package digital

import (
	"context"
	"ss-catalog-service/internal/domain"
)

type digitalUsecase struct {
	repo domain.DigitalRepository
}

func NewDigitalUsecase(repo domain.DigitalRepository) domain.DigitalUsecase {
	return &digitalUsecase{repo: repo}
}

func (u *digitalUsecase) UploadDigitalProduct(ctx context.Context, file *domain.DigitalFile) error {
	return u.repo.AddFile(ctx, file)
}

func (u *digitalUsecase) AddLicenses(ctx context.Context, productID int, keys []string) error {
	licenseKeys := make([]domain.LicenseKey, len(keys))
	for i, key := range keys {
		licenseKeys[i] = domain.LicenseKey{
			ProductID:  productID,
			LicenseKey: key,
			IsSold:     false,
		}
	}
	return u.repo.AddLicenseKeys(ctx, licenseKeys)
}

func (u *digitalUsecase) GetDigitalDetails(ctx context.Context, productID int) ([]domain.DigitalFile, int, error) {
	files, err := u.repo.GetFilesByProductID(ctx, productID)
	if err != nil {
		return nil, 0, err
	}

	count, err := u.repo.GetAvailableLicenseCount(ctx, productID)
	if err != nil {
		return files, 0, nil // Return files even if license count fails
	}

	return files, count, nil
}

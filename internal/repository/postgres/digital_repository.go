package postgres

import (
	"context"
	"ss-catalog-service/internal/domain"
	"time"

	"gorm.io/gorm"
)

type digitalRepository struct {
	db *gorm.DB
}

func NewDigitalRepository(db *gorm.DB) domain.DigitalRepository {
	return &digitalRepository{db: db}
}

func (r *digitalRepository) AddFile(ctx context.Context, file *domain.DigitalFile) error {
	model := &DigitalFileModel{
		ProductID:     file.ProductID,
		FileName:      file.FileName,
		FilePath:      file.FilePath,
		FileSizeBytes: file.FileSizeBytes,
		MimeType:      file.MimeType,
		Version:       file.Version,
	}
	db := getDB(ctx, r.db)
	return db.Create(model).Error
}

func (r *digitalRepository) GetFilesByProductID(ctx context.Context, productID int) ([]domain.DigitalFile, error) {
	var models []DigitalFileModel
	db := getDB(ctx, r.db)
	if err := db.Where("product_id = ?", productID).Find(&models).Error; err != nil {
		return nil, err
	}
	files := make([]domain.DigitalFile, len(models))
	for i, m := range models {
		files[i] = m.ToDomain()
	}
	return files, nil
}

func (r *digitalRepository) AddLicenseKeys(ctx context.Context, keys []domain.LicenseKey) error {
	models := make([]LicenseKeyModel, len(keys))
	for i, k := range keys {
		models[i] = LicenseKeyModel{
			ProductID:  k.ProductID,
			LicenseKey: k.LicenseKey,
		}
	}
	db := getDB(ctx, r.db)
	return db.Create(&models).Error
}

func (r *digitalRepository) GetAvailableLicenseCount(ctx context.Context, productID int) (int, error) {
	var count int64
	db := getDB(ctx, r.db)
	err := db.Model(&LicenseKeyModel{}).Where("product_id = ? AND is_sold = false", productID).Count(&count).Error
	return int(count), err
}

func (r *digitalRepository) AssignLicenseKey(ctx context.Context, productID int, orderID string) (*domain.LicenseKey, error) {
	var model LicenseKeyModel
	db := getDB(ctx, r.db)
	
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ? AND is_sold = false", productID).First(&model).Error; err != nil {
			return err
		}
		now := time.Now()
		updates := map[string]interface{}{
			"is_sold":  true,
			"sold_at":  &now,
			"order_id": orderID,
		}
		return tx.Model(&model).Updates(updates).Error
	})

	if err != nil {
		return nil, err
	}
	key := model.ToDomain()
	return &key, nil
}

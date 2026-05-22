package postgres

import (
	"context"
	"log/slog"
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
	err := db.Create(model).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "add_digital_file", "product_id", file.ProductID, "error", err)
		return err
	}
	file.ID = model.ID
	file.CreatedAt = model.CreatedAt
	slog.InfoContext(ctx, "db_audit", "operation", "add_digital_file", "file_id", file.ID, "product_id", file.ProductID)
	return nil
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
	err := db.Create(&models).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "add_license_keys", "error", err)
		return err
	}
	slog.InfoContext(ctx, "db_audit", "operation", "add_license_keys", "count", len(keys))
	return nil
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
		slog.ErrorContext(ctx, "db_error", "operation", "assign_license_key", "product_id", productID, "order_id", orderID, "error", err)
		return nil, err
	}
	key := model.ToDomain()
	slog.InfoContext(ctx, "db_audit", "operation", "assign_license_key", "license_key_id", key.ID, "product_id", productID, "order_id", orderID)
	return &key, nil
}

package postgres

import (
	"ss-catalog-service/internal/domain"
	"time"
)

type DigitalFileModel struct {
	ID             int       `gorm:"primaryKey"`
	ProductID      int       `gorm:"index"`
	FileName       string    `gorm:"type:varchar(255);not null"`
	FilePath       string    `gorm:"type:varchar(500);not null"`
	FileSizeBytes  int64     `gorm:"not null"`
	MimeType       string    `gorm:"type:varchar(100)"`
	Version        string    `gorm:"type:varchar(50)"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (DigitalFileModel) TableName() string {
	return "product_digital_files"
}

type LicenseKeyModel struct {
	ID         int        `gorm:"primaryKey"`
	ProductID  int        `gorm:"index"`
	LicenseKey string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	IsSold     bool       `gorm:"default:false;index"`
	SoldAt     *time.Time
	OrderID    string     `gorm:"type:varchar(100);index"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
}

func (LicenseKeyModel) TableName() string {
	return "product_license_keys"
}

func (m *DigitalFileModel) ToDomain() domain.DigitalFile {
	return domain.DigitalFile{
		ID:            m.ID,
		ProductID:     m.ProductID,
		FileName:      m.FileName,
		FilePath:      m.FilePath,
		FileSizeBytes: m.FileSizeBytes,
		MimeType:      m.MimeType,
		Version:       m.Version,
		CreatedAt:     m.CreatedAt,
	}
}

func (m *LicenseKeyModel) ToDomain() domain.LicenseKey {
	return domain.LicenseKey{
		ID:         m.ID,
		ProductID:  m.ProductID,
		LicenseKey: m.LicenseKey,
		IsSold:     m.IsSold,
		SoldAt:     m.SoldAt,
		OrderID:    m.OrderID,
		CreatedAt:  m.CreatedAt,
	}
}

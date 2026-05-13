package domain

import (
	"context"
	"time"
)

type DigitalFile struct {
	ID             int
	ProductID      int
	FileName       string
	FilePath       string
	FileSizeBytes  int64
	MimeType       string
	Version        string
	CreatedAt      time.Time
}

type LicenseKey struct {
	ID        int
	ProductID int
	LicenseKey string
	IsSold    bool
	SoldAt    *time.Time
	OrderID   string
	CreatedAt time.Time
}

// --- Interfaces ---

type DigitalRepository interface {
	AddFile(ctx context.Context, file *DigitalFile) error
	GetFilesByProductID(ctx context.Context, productID int) ([]DigitalFile, error)
	AddLicenseKeys(ctx context.Context, keys []LicenseKey) error
	GetAvailableLicenseCount(ctx context.Context, productID int) (int, error)
	AssignLicenseKey(ctx context.Context, productID int, orderID string) (*LicenseKey, error)
}

type DigitalUsecase interface {
	UploadDigitalProduct(ctx context.Context, file *DigitalFile) error
	AddLicenses(ctx context.Context, productID int, keys []string) error
	GetDigitalDetails(ctx context.Context, productID int) ([]DigitalFile, int, error)
}

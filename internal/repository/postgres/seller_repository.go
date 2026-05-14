package postgres

import (
	"context"
	"errors"
	"ss-catalog-service/internal/domain"
	"time"

	"gorm.io/gorm"
)

type SellerModel struct {
	BaseModel
	Name       string     `gorm:"type:varchar(255);not null"`
	Code       string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	IsActive   bool       `gorm:"not null;default:true"`
	VerifiedAt *time.Time `gorm:"null"`
}

func (SellerModel) TableName() string { return "sellers" }

type SellerUserModel struct {
	BaseModel
	SellerID int    `gorm:"not null;index"`
	UserID   int    `gorm:"not null;index"`
	Role     string `gorm:"type:varchar(50);not null;default:'owner'"`
}

func (SellerUserModel) TableName() string { return "seller_users" }

type SellerProductModel struct {
	BaseModel
	SellerID   int        `gorm:"not null;index"`
	ProductID  int        `gorm:"not null;index"`
	IsActive   bool       `gorm:"not null;default:true"`
	ApprovedAt *time.Time `gorm:"null"`
	ApprovedBy string     `gorm:"type:varchar(255)"`
}

func (SellerProductModel) TableName() string { return "seller_products" }

type sellerRepository struct {
	db *gorm.DB
}

func NewSellerRepository(db *gorm.DB) domain.SellerRepository {
	return &sellerRepository{db: db}
}

func (r *sellerRepository) FindSellerIDByUserID(ctx context.Context, userID int) (int, error) {
	var mapping SellerUserModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		First(&mapping).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // No seller associated
		}
		return 0, err
	}

	return mapping.SellerID, nil
}

func (r *sellerRepository) FindByID(ctx context.Context, id int) (*domain.Seller, error) {
	var m SellerModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}

	return &domain.Seller{
		BaseEntity: domain.BaseEntity{
			ID:        m.ID,
			PublicID:  m.PublicID,
			CreatedAt: m.CreatedAt,
		},
		Name:       m.Name,
		Code:       m.Code,
		IsActive:   m.IsActive,
		VerifiedAt: m.VerifiedAt,
	}, nil
}

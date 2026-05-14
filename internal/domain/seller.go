package domain

import (
	"context"
	"time"
)

// Seller represents the business entity for a vendor/store.
type Seller struct {
	BaseEntity
	Name       string     `json:"name"`
	Code       string     `json:"code"`
	IsActive   bool       `json:"is_active"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// SellerUser represents the mapping between a User and a Seller.
type SellerUser struct {
	BaseEntity
	SellerID int    `json:"seller_id"`
	UserID   int    `json:"user_id"`
	Role     string `json:"role"`
}

// SellerProduct maps which sellers are authorized to sell each product.
type SellerProduct struct {
	BaseEntity
	SellerID   int        `json:"seller_id"`
	ProductID  int        `json:"product_id"`
	IsActive   bool       `json:"is_active"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	ApprovedBy string     `json:"approved_by,omitempty"`
}

// SellerRepository defines the contract for seller-related data access.
type SellerRepository interface {
	FindSellerIDByUserID(ctx context.Context, userID int) (int, error)
	FindByID(ctx context.Context, id int) (*Seller, error)
}

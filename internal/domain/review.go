package domain

import (
	"context"
	"time"
)

type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
)

type ProductReview struct {
	BaseEntity
	ProductID          int
	UserID             string
	UserName           string
	Rating             int
	Comment            string
	IsVerifiedPurchase bool
	Status             ReviewStatus
	Images             []ReviewImage
	Votes              []ReviewVote
}

type ReviewImage struct {
	ID        int
	ReviewID  int
	ImageURL  string
	CreatedAt time.Time
}

type ReviewVote struct {
	ReviewID  int
	UserID    string
	IsHelpful bool
	CreatedAt time.Time
}

// --- Interfaces ---

type ReviewRepository interface {
	Create(ctx context.Context, review *ProductReview) error
	GetByProductID(ctx context.Context, productID int, p Pagination) ([]ProductReview, error)
	GetAverageRating(ctx context.Context, productID int) (float64, int, error)
	AddVote(ctx context.Context, vote *ReviewVote) error
	UpdateStatus(ctx context.Context, reviewID int, status ReviewStatus) error
	FindAll(ctx context.Context, p Pagination) ([]ProductReview, int64, error)
}

type ReviewUsecase interface {
	SubmitReview(ctx context.Context, review *ProductReview) error
	GetProductReviews(ctx context.Context, productID int, p Pagination) ([]ProductReview, error)
	GetProductRatingSummary(ctx context.Context, productID int) (float64, int, error)
	VoteReview(ctx context.Context, reviewID int, userID string, helpful bool) error
	GetAllReviews(ctx context.Context, p Pagination) ([]ProductReview, int64, error)
	UpdateReviewStatus(ctx context.Context, reviewID int, status ReviewStatus) error
}

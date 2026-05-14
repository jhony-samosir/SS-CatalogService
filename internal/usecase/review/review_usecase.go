package review

import (
	"context"
	"ss-catalog-service/internal/domain"
)

type reviewUsecase struct {
	repo domain.ReviewRepository
}

func NewReviewUsecase(repo domain.ReviewRepository) domain.ReviewUsecase {
	return &reviewUsecase{repo: repo}
}

func (u *reviewUsecase) SubmitReview(ctx context.Context, review *domain.ProductReview) error {
	// Add business validation here (e.g. check if user already reviewed)
	review.Status = domain.ReviewStatusPending // Default status
	return u.repo.Create(ctx, review)
}

func (u *reviewUsecase) GetProductReviews(ctx context.Context, productID int, p domain.Pagination) ([]domain.ProductReview, error) {
	return u.repo.GetByProductID(ctx, productID, p)
}

func (u *reviewUsecase) GetProductRatingSummary(ctx context.Context, productID int) (float64, int, error) {
	return u.repo.GetAverageRating(ctx, productID)
}

func (u *reviewUsecase) VoteReview(ctx context.Context, reviewID int, userID string, helpful bool) error {
	vote := &domain.ReviewVote{
		ReviewID:  reviewID,
		UserID:    userID,
		IsHelpful: helpful,
	}
	return u.repo.AddVote(ctx, vote)
}

func (u *reviewUsecase) GetAllReviews(ctx context.Context, p domain.Pagination) ([]domain.ProductReview, int64, error) {
	return u.repo.FindAll(ctx, p)
}

func (u *reviewUsecase) UpdateReviewStatus(ctx context.Context, reviewID int, status domain.ReviewStatus) error {
	return u.repo.UpdateStatus(ctx, reviewID, status)
}

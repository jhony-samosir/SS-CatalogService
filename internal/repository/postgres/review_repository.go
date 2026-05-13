package postgres

import (
	"context"
	"ss-catalog-service/internal/domain"

	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) domain.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *domain.ProductReview) error {
	model := FromReviewDomain(review)
	db := getDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		return mapDBError(err)
	}

	review.ID = model.ID
	review.CreatedAt = model.CreatedAt
	return nil
}

func (r *reviewRepository) GetByProductID(ctx context.Context, productID int, p domain.Pagination) ([]domain.ProductReview, error) {
	var models []ProductReviewModel
	db := getDB(ctx, r.db)

	query := db.Where("product_id = ? AND status = ?", productID, domain.ReviewStatusApproved).
		Preload("Images").
		Preload("Votes")

	if p.Limit > 0 {
		query = query.Limit(p.Limit).Offset(p.Offset)
	}

	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	reviews := make([]domain.ProductReview, len(models))
	for i, m := range models {
		reviews[i] = m.ToDomain()
	}
	return reviews, nil
}

func (r *reviewRepository) GetAverageRating(ctx context.Context, productID int) (float64, int, error) {
	db := getDB(ctx, r.db)
	
	var result struct {
		AvgRating float64
		Count     int
	}

	err := db.Model(&ProductReviewModel{}).
		Select("AVG(rating) as avg_rating, COUNT(*) as count").
		Where("product_id = ? AND status = ?", productID, domain.ReviewStatusApproved).
		Scan(&result).Error

	if err != nil {
		return 0, 0, err
	}

	return result.AvgRating, result.Count, nil
}

func (r *reviewRepository) AddVote(ctx context.Context, vote *domain.ReviewVote) error {
	model := &ReviewVoteModel{
		ReviewID:  vote.ReviewID,
		UserID:    vote.UserID,
		IsHelpful: vote.IsHelpful,
	}
	db := getDB(ctx, r.db)

	return db.Save(model).Error
}

func (r *reviewRepository) UpdateStatus(ctx context.Context, reviewID int, status domain.ReviewStatus) error {
	db := getDB(ctx, r.db)
	return db.Model(&ProductReviewModel{}).Where("id = ?", reviewID).Update("status", string(status)).Error
}

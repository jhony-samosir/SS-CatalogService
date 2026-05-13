package postgres

import (
	"ss-catalog-service/internal/domain"
	"time"

	"gorm.io/gorm"
)

type ProductReviewModel struct {
	BaseModel
	ProductID          int           `gorm:"index:idx_product_reviews_product_id"`
	UserID             string        `gorm:"type:varchar(255);index:idx_product_reviews_user_id"`
	UserName           string        `gorm:"type:varchar(100)"`
	Rating             int           `gorm:"type:smallint;check:rating >= 1 AND rating <= 5"`
	Comment            string        `gorm:"type:text"`
	IsVerifiedPurchase bool          `gorm:"default:false"`
	Status             string        `gorm:"type:varchar(20);default:'pending'"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
	DeletedBy          string        `gorm:"type:varchar(255)"`
	
	Images []ReviewImageModel `gorm:"foreignKey:ReviewID"`
	Votes  []ReviewVoteModel  `gorm:"foreignKey:ReviewID"`
}

func (ProductReviewModel) TableName() string {
	return "product_reviews"
}

type ReviewImageModel struct {
	ID        int       `gorm:"primaryKey"`
	ReviewID  int       `gorm:"index"`
	ImageURL  string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (ReviewImageModel) TableName() string {
	return "product_review_images"
}

type ReviewVoteModel struct {
	ReviewID  int       `gorm:"primaryKey;autoIncrement:false"`
	UserID    string    `gorm:"primaryKey;type:varchar(255)"`
	IsHelpful bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (ReviewVoteModel) TableName() string {
	return "product_review_votes"
}

func (m *ProductReviewModel) ToDomain() domain.ProductReview {
	review := domain.ProductReview{
		BaseEntity: domain.BaseEntity{
			ID:        m.ID,
			PublicID:  m.PublicID,
			CreatedAt: m.CreatedAt,
			CreatedBy: m.CreatedBy,
			UpdatedAt: m.UpdatedAt,
			UpdatedBy: m.UpdatedBy,
		},
		ProductID:          m.ProductID,
		UserID:             m.UserID,
		UserName:           m.UserName,
		Rating:             m.Rating,
		Comment:            m.Comment,
		IsVerifiedPurchase: m.IsVerifiedPurchase,
		Status:             domain.ReviewStatus(m.Status),
	}

	if m.DeletedAt.Valid {
		review.DeletedAt = &m.DeletedAt.Time
		review.DeletedBy = m.DeletedBy
	}

	for _, img := range m.Images {
		review.Images = append(review.Images, domain.ReviewImage{
			ID:        img.ID,
			ReviewID:  img.ReviewID,
			ImageURL:  img.ImageURL,
			CreatedAt: img.CreatedAt,
		})
	}

	for _, vote := range m.Votes {
		review.Votes = append(review.Votes, domain.ReviewVote{
			ReviewID:  vote.ReviewID,
			UserID:    vote.UserID,
			IsHelpful: vote.IsHelpful,
			CreatedAt: vote.CreatedAt,
		})
	}

	return review
}

func FromReviewDomain(d *domain.ProductReview) *ProductReviewModel {
	m := &ProductReviewModel{
		BaseModel: BaseModel{
			ID:        d.ID,
			PublicID:  d.PublicID,
			CreatedAt: d.CreatedAt,
			CreatedBy: d.CreatedBy,
			UpdatedBy: d.UpdatedBy,
		},
		ProductID:          d.ProductID,
		UserID:             d.UserID,
		UserName:           d.UserName,
		Rating:             d.Rating,
		Comment:            d.Comment,
		IsVerifiedPurchase: d.IsVerifiedPurchase,
		Status:             string(d.Status),
	}

	if d.UpdatedAt != nil {
		m.UpdatedAt = d.UpdatedAt
	}

	if d.DeletedAt != nil {
		m.DeletedAt = gorm.DeletedAt{Time: *d.DeletedAt, Valid: true}
		m.DeletedBy = d.DeletedBy
	}

	return m
}

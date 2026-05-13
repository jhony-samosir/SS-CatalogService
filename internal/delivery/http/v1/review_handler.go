package v1

import (
	"net/http"
	"strconv"
	"ss-catalog-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	usecase        domain.ReviewUsecase
	productQueries domain.ProductQueryUsecase
}

func NewReviewHandler(u domain.ReviewUsecase, pq domain.ProductQueryUsecase) *ReviewHandler {
	return &ReviewHandler{
		usecase:        u,
		productQueries: pq,
	}
}

type SubmitReviewRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"required,max=1000"`
	UserName  string `json:"user_name" binding:"required"`
}

func (h *ReviewHandler) SubmitReview(c *gin.Context) {
	var req SubmitReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publicID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product uuid format"})
		return
	}

	product, err := h.productQueries.GetProductByPublicID(c.Request.Context(), publicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// In a real app, UserID would come from JWT claims
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "anonymous" // Placeholder for testing
	}

	err = h.usecase.SubmitReview(c.Request.Context(), &domain.ProductReview{
		ProductID: product.ID,
		UserID:    userID,
		UserName:  req.UserName,
		Rating:    req.Rating,
		Comment:   req.Comment,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "review submitted successfully"})
}

func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	publicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product uuid format"})
		return
	}

	product, err := h.productQueries.GetProductByPublicID(c.Request.Context(), publicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	productID := product.ID

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reviews, err := h.usecase.GetProductReviews(c.Request.Context(), productID, domain.Pagination{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reviews})
}

func (h *ReviewHandler) GetRatingSummary(c *gin.Context) {
	publicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product uuid format"})
		return
	}

	product, err := h.productQueries.GetProductByPublicID(c.Request.Context(), publicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	productID := product.ID

	avg, count, err := h.usecase.GetProductRatingSummary(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get rating summary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"average_rating": avg,
		"total_reviews":  count,
	})
}

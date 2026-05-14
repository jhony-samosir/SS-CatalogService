package v1

import (
	"net/http"
	"strconv"
	"ss-catalog-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BundleHandler struct {
	usecase domain.BundleUsecase
}

func NewBundleHandler(u domain.BundleUsecase) *BundleHandler {
	return &BundleHandler{usecase: u}
}

type CreateBundleRequest struct {
	Name          string       `json:"name" binding:"required"`
	Slug          string       `json:"slug" binding:"required"`
	Description   string       `json:"description"`
	PriceOverride *float64     `json:"price_override"`
	Items         []BundleItem `json:"items" binding:"required,min=1"`
}

type BundleItem struct {
	ProductID *int `json:"product_id"`
	VariantID *int `json:"variant_id"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func (h *BundleHandler) CreateBundle(c *gin.Context) {
	var req CreateBundleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items := make([]domain.BundleItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.BundleItem{
			ProductID: item.ProductID,
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		}
	}

	err := h.usecase.CreateBundle(c.Request.Context(), &domain.ProductBundle{
		Name:          req.Name,
		Slug:          req.Slug,
		Description:   req.Description,
		PriceOverride: req.PriceOverride,
		IsActive:      true,
		Items:         items,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create bundle"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "bundle created successfully"})
}

func (h *BundleHandler) GetBundles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	bundles, total, err := h.usecase.GetBundles(c.Request.Context(), domain.Pagination{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get bundles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"items":       bundles,
			"total_count": total,
		},
	})
}

func (h *BundleHandler) GetBundle(c *gin.Context) {
	publicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bundle id"})
		return
	}

	bundle, err := h.usecase.GetBundleByPublicID(c.Request.Context(), publicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "bundle not found"})
		return
	}

	c.JSON(http.StatusOK, bundle)
}

func (h *BundleHandler) UpdateBundle(c *gin.Context) {
	publicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bundle id"})
		return
	}

	var req CreateBundleRequest // Use same structure for simplicity or a separate UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items := make([]domain.BundleItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.BundleItem{
			ProductID: item.ProductID,
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		}
	}

	err = h.usecase.UpdateBundle(c.Request.Context(), &domain.ProductBundle{
		PublicID:      publicID,
		Name:          req.Name,
		Slug:          req.Slug,
		Description:   req.Description,
		PriceOverride: req.PriceOverride,
		Items:         items,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update bundle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bundle updated successfully"})
}

func (h *BundleHandler) DeleteBundle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.usecase.DeleteBundle(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete bundle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bundle deleted successfully"})
}

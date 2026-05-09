package v1

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"ss-catalog-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	commandUsecase domain.ProductCommandUsecase
	queryUsecase   domain.ProductQueryUsecase
}

func NewProductHandler(cu domain.ProductCommandUsecase, qu domain.ProductQueryUsecase) *ProductHandler {
	return &ProductHandler{
		commandUsecase: cu,
		queryUsecase:   qu,
	}
}

// --- Request / Response DTOs ---

type CreateProductRequest struct {
	Name    string `json:"name" binding:"required,min=3,max=500"`
	BrandID *int   `json:"brand_id,omitempty"`
}

type UpdateProductRequest struct {
	Name        string               `json:"name" binding:"required,min=3,max=500"`
	Description string               `json:"description"`
	Status      domain.ProductStatus `json:"status" binding:"required"`
}

type ProductResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Slug    string    `json:"slug"`
	Status  string    `json:"status"`
	BrandID *int      `json:"brand_id,omitempty"`
}

func toProductResponse(p *domain.Product) ProductResponse {
	return ProductResponse{
		ID:      p.PublicID,
		Name:    p.Name,
		Slug:    p.Slug,
		Status:  string(p.Status),
		BrandID: p.BrandID,
	}
}

// CreateProduct handles POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.commandUsecase.CreateProduct(c.Request.Context(), domain.CreateProductPayload{
		Name:    req.Name,
		BrandID: req.BrandID,
	})
	if err != nil {
		log.Printf("ERROR [CreateProduct]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(product))
}

// UpdateProduct handles PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	publicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id format"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.commandUsecase.UpdateProduct(c.Request.Context(), domain.UpdateProductPayload{
		PublicID:    publicID,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	})

	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if errors.Is(err, domain.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to update this product"})
			return
		}
		log.Printf("ERROR [UpdateProduct]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
}

// GetProduct handles GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	publicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id format"})
		return
	}

	product, err := h.queryUsecase.GetProductByPublicID(c.Request.Context(), publicID)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		log.Printf("ERROR [GetProduct]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve product"})
		return
	}

	// Add Caching Headers
	etag := fmt.Sprintf("\"%d\"", product.UpdatedAt.Unix())
	c.Header("Cache-Control", "public, max-age=60")
	c.Header("Vary", "Accept-Language, Authorization")
	
	if product.UpdatedAt != nil {
		c.Header("ETag", etag)
		// 304 Not Modified support
		if c.GetHeader("If-None-Match") == etag {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}
	}

	c.JSON(http.StatusOK, toProductResponse(product))
}

// GetProducts handles GET /api/v1/products?limit=10&offset=0
func (h *ProductHandler) GetProducts(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	// Hard cap for safety
	if limit > 100 {
		limit = 100
	}

	products, err := h.queryUsecase.GetAllProducts(c.Request.Context(), domain.Pagination{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Printf("ERROR [GetProducts]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve products"})
		return
	}

	response := make([]ProductResponse, len(products))
	for i := range products {
		response[i] = toProductResponse(&products[i])
	}

	// Add Caching Headers
	c.Header("Cache-Control", "public, max-age=60")
	c.Header("Vary", "Accept-Language, Authorization")

	c.JSON(http.StatusOK, gin.H{
		"data":   response,
		"limit":  limit,
		"offset": offset,
	})
}

// SearchProducts handles GET /api/v1/products/search
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	q := domain.GetProductSearchQuery{}

	if kw := c.Query("q"); kw != "" {
		q.Keyword = &kw
	}
	if cs := c.Query("category_slug"); cs != "" {
		q.CategorySlug = &cs
	}
	if bid := c.Query("brand_id"); bid != "" {
		id, err := strconv.Atoi(bid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "brand_id must be a valid integer"})
			return
		}
		q.BrandID = &id
	}
	if minP := c.Query("min_price"); minP != "" {
		v, err := strconv.ParseFloat(minP, 64)
		if err != nil || v < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_price"})
			return
		}
		q.MinPrice = &v
	}
	if maxP := c.Query("max_price"); maxP != "" {
		v, err := strconv.ParseFloat(maxP, 64)
		if err != nil || v < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_price"})
			return
		}
		q.MaxPrice = &v
	}
	if cur := c.Query("cursor"); cur != "" {
		q.Cursor = &cur
	}
	if lim, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		q.Limit = lim
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	result, err := h.queryUsecase.SearchProducts(ctx, q)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCursor) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired pagination cursor"})
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("ERROR [SearchProducts]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	items := make([]ProductResponse, len(result.Items))
	for i := range result.Items {
		items[i] = toProductResponse(&result.Items[i])
	}

	// Add Caching Headers
	c.Header("Cache-Control", "public, max-age=30") // Shorter for search results
	c.Header("Vary", "Accept-Language, Authorization")

	c.JSON(http.StatusOK, gin.H{
		"data":        items,
		"next_cursor": result.NextCursor,
		"has_more":    result.NextCursor != nil,
	})
}

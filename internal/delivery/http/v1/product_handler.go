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
	"ss-catalog-service/pkg/response"

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
	Name        string   `json:"name" binding:"required,min=3,max=500"`
	Slug        string   `json:"slug"`
	Price       float64  `json:"price"`
	ImageURL    string   `json:"image_url"`
	Description string   `json:"description"`
	BrandID     *string  `json:"brand_id,omitempty"`
	Status      string   `json:"status"`
	CategoryIDs []string `json:"category_ids"`
}

type UpdateProductRequest struct {
	Name        string               `json:"name" binding:"required,min=3,max=500"`
	Description string               `json:"description"`
	Status      domain.ProductStatus `json:"status" binding:"required"`
	BrandID     *string              `json:"brand_id,omitempty"`
	CategoryIDs []string             `json:"category_ids"`
}

type ProductResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Status   string    `json:"status"`
	Price    float64   `json:"price"`
	ImageURL string    `json:"image_url"`
	Rating   float64   `json:"rating"`
	BrandID  *int      `json:"brand_id,omitempty"`
}

func toProductResponse(p *domain.Product) ProductResponse {
	return ProductResponse{
		ID:       p.PublicID,
		Name:     p.Name,
		Slug:     p.Slug,
		Status:   string(p.Status),
		Price:    0, // SPU level usually doesn't have price, variants do. Defaulting to 0.
		ImageURL: p.ImageURL,
		Rating:   0, // Defaulting to 0
		BrandID:  p.BrandID,
	}
}

// CreateProduct handles POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := domain.CreateProductPayload{
		Name:        req.Name,
		Slug:        req.Slug,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
		Description: req.Description,
		Status:      domain.ProductStatus(req.Status),
	}

	if req.BrandID != nil && *req.BrandID != "" {
		bid, err := uuid.Parse(*req.BrandID)
		if err == nil {
			payload.PublicBrandID = &bid
		}
	}

	if len(req.CategoryIDs) > 0 {
		payload.CategoryPublicIDs = make([]uuid.UUID, 0, len(req.CategoryIDs))
		for _, idStr := range req.CategoryIDs {
			if id, err := uuid.Parse(idStr); err == nil {
				payload.CategoryPublicIDs = append(payload.CategoryPublicIDs, id)
			}
		}
	}

	product, err := h.commandUsecase.CreateProduct(c.Request.Context(), payload)
	if err != nil {
		log.Printf("ERROR [CreateProduct]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product: " + err.Error()})
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

	payload := domain.UpdateProductPayload{
		PublicID:    publicID,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	if req.BrandID != nil && *req.BrandID != "" {
		bid, err := uuid.Parse(*req.BrandID)
		if err == nil {
			payload.PublicBrandID = &bid
		}
	}

	if len(req.CategoryIDs) > 0 {
		payload.CategoryPublicIDs = make([]uuid.UUID, 0, len(req.CategoryIDs))
		for _, idStr := range req.CategoryIDs {
			if id, err := uuid.Parse(idStr); err == nil {
				payload.CategoryPublicIDs = append(payload.CategoryPublicIDs, id)
			}
		}
	}

	err = h.commandUsecase.UpdateProduct(c.Request.Context(), payload)

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

// GetProducts handles GET /api/v1/products?page=1&limit=10
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var p domain.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination parameters"})
		return
	}
	p.SetDefaults()

	products, total, err := h.queryUsecase.GetAllProducts(c.Request.Context(), p)
	if err != nil {
		log.Printf("ERROR [GetProducts]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve products"})
		return
	}

	responseItems := make([]ProductResponse, len(products))
	for i := range products {
		responseItems[i] = toProductResponse(&products[i])
	}

	// Add Caching Headers
	c.Header("Cache-Control", "public, max-age=60")
	c.Header("Vary", "Accept-Language, Authorization")

	response.PaginatedJSON(c, http.StatusOK, "Products retrieved successfully", responseItems, total, p.Page, p.Limit)
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

// FacetedSearch handles GET /api/v1/products/faceted-search
func (h *ProductHandler) FacetedSearch(c *gin.Context) {
	q := domain.GetProductSearchQuery{}
	// (Reuse logic from SearchProducts or refactor)
	if kw := c.Query("q"); kw != "" {
		q.Keyword = &kw
	}
	if lim, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		q.Limit = lim
	}

	result, err := h.queryUsecase.FacetedSearch(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "faceted search failed"})
		return
	}

	c.JSON(http.StatusOK, result)
}

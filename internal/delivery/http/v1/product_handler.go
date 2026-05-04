package v1

import (
	"errors"
	"log"
	"net/http"
	"strconv"
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

	c.JSON(http.StatusOK, gin.H{
		"data":   response,
		"limit":  limit,
		"offset": offset,
	})
}

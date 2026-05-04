package v1

import (
	"errors"
	"log"
	"net/http"
	"ss-catalog-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VariantHandler struct {
	commandUsecase domain.VariantCommandUsecase
}

func NewVariantHandler(cu domain.VariantCommandUsecase) *VariantHandler {
	return &VariantHandler{
		commandUsecase: cu,
	}
}

// --- Request DTOs ---

type CreateVariantRequest struct {
	ProductID  int                        `json:"product_id" binding:"required"`
	SKU        string                     `json:"sku" binding:"required,min=3,max=100"`
	Barcode    string                     `json:"barcode"`
	Name       string                     `json:"name"`
	IsDefault  bool                       `json:"is_default"`
	WeightGram *int                       `json:"weight_gram"`
	Attributes []VariantAttributeRequest `json:"attributes"`
	Images     []VariantImageRequest     `json:"images"`
}

type VariantAttributeRequest struct {
	AttributeID      int `json:"attribute_id" binding:"required"`
	AttributeValueID int `json:"attribute_value_id" binding:"required"`
}

type VariantImageRequest struct {
	URL       string `json:"url" binding:"required,url"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary"`
}

// --- Response DTOs ---

type VariantResponse struct {
	ID        uuid.UUID `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	IsDefault bool      `json:"is_default"`
}

func toVariantResponse(v *domain.ProductVariant) VariantResponse {
	return VariantResponse{
		ID:        v.PublicID,
		SKU:       v.SKU,
		Name:      v.Name,
		IsDefault: v.IsDefault,
	}
}

// CreateVariant handles POST /api/v1/variants
func (h *VariantHandler) CreateVariant(c *gin.Context) {
	var req CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map Request to Domain Payload
	payload := domain.CreateVariantPayload{
		ProductID:  req.ProductID,
		SKU:        req.SKU,
		Barcode:    req.Barcode,
		Name:       req.Name,
		IsDefault:  req.IsDefault,
		WeightGram: req.WeightGram,
	}

	for _, attr := range req.Attributes {
		payload.Attributes = append(payload.Attributes, domain.VariantAttributePayload{
			AttributeID:      attr.AttributeID,
			AttributeValueID: attr.AttributeValueID,
		})
	}

	for _, img := range req.Images {
		payload.Images = append(payload.Images, domain.VariantImagePayload{
			URL:       img.URL,
			AltText:   img.AltText,
			IsPrimary: img.IsPrimary,
		})
	}

	variant, err := h.commandUsecase.CreateProductVariant(c.Request.Context(), payload)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateSKU) {
			c.JSON(http.StatusConflict, gin.H{"error": "SKU already exists"})
			return
		}
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "parent product (SPU) not found"})
			return
		}

		log.Printf("ERROR [CreateVariant]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create variant"})
		return
	}

	c.JSON(http.StatusCreated, toVariantResponse(variant))
}

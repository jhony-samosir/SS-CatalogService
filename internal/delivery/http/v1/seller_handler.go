package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SellerHandler struct {
	// Usecases to be added later
}

func NewSellerHandler() *SellerHandler {
	return &SellerHandler{}
}

func (h *SellerHandler) RegisterSeller(c *gin.Context) {
	// Placeholder for future implementation
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Seller registration not yet implemented"})
}

func (h *SellerHandler) GetSeller(c *gin.Context) {
	// Placeholder for future implementation
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Get seller not yet implemented"})
}

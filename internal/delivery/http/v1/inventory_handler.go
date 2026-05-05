package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ss-catalog-service/internal/domain"
)

type InventoryHandler struct {
	cmdUsecase domain.InventoryCommandUsecase
}

func NewInventoryHandler(cmdUsecase domain.InventoryCommandUsecase) *InventoryHandler {
	return &InventoryHandler{cmdUsecase: cmdUsecase}
}

func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	var payload domain.UpdateStockPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}

	if payload.VariantID <= 0 || payload.WarehouseID <= 0 || payload.Quantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "variant_id, warehouse_id, and non-zero quantity are required"})
		return
	}

	if err := h.cmdUsecase.UpdateInventoryStock(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock adjusted successfully"})
}

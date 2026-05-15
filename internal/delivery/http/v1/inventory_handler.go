package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ss-catalog-service/internal/domain"
	"ss-catalog-service/pkg/response"
)

type InventoryHandler struct {
	cmdUsecase domain.InventoryCommandUsecase
	queryUsecase domain.InventoryQueryUsecase
}

func NewInventoryHandler(cmdUsecase domain.InventoryCommandUsecase, queryUsecase domain.InventoryQueryUsecase) *InventoryHandler {
	return &InventoryHandler{
		cmdUsecase: cmdUsecase,
		queryUsecase: queryUsecase,
	}
}

func (h *InventoryHandler) GetInventory(c *gin.Context) {
	var p domain.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p.SetDefaults()
	}
	p.SetDefaults()

	warehouseID := c.Query("warehouse_id")
	variantID := c.Query("variant_id")

	items, total, err := h.queryUsecase.GetInventory(c.Request.Context(), p, warehouseID, variantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch inventory", err.Error())
		return
	}

	response.PaginatedJSON(c, http.StatusOK, "Inventory fetched successfully", items, total, p.Page, p.Limit)
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

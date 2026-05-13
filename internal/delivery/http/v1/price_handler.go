package v1

import (
	"net/http"
	"strconv"
	"ss-catalog-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type PriceHandler struct {
	repo domain.PriceHistoryRepository
}

func NewPriceHandler(r domain.PriceHistoryRepository) *PriceHandler {
	return &PriceHandler{repo: r}
}

func (h *PriceHandler) GetPriceHistory(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var variantID *int
	if vID := c.Query("variant_id"); vID != "" {
		id, err := strconv.Atoi(vID)
		if err == nil {
			variantID = &id
		}
	}

	history, err := h.repo.GetPriceHistory(c.Request.Context(), productID, variantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get price history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}

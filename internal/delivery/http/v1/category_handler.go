package v1

import (
	"net/http"
	"ss-catalog-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	usecase domain.CategoryUsecase
}

func NewCategoryHandler(u domain.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: u}
}

// GetCategories returns a list of all categories.
// @Summary Get all categories
// @Description Get a list of all product categories with pagination
// @Tags Categories
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {array} domain.Category
// @Router /api/catalog/v1/categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var p domain.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p = domain.Pagination{Limit: 100, Offset: 0}
	}

	categories, err := h.usecase.GetCategories(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

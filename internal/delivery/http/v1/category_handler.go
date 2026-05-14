package v1

import (
	"net/http"
	"ss-catalog-service/internal/delivery/http/v1/dto"
	"ss-catalog-service/internal/domain"
	"ss-catalog-service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	usecase domain.CategoryUsecase
}

func NewCategoryHandler(u domain.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: u}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var p domain.Pagination
	_ = c.ShouldBindQuery(&p)
	p.SetDefaults()

	categories, total, err := h.usecase.GetCategories(c.Request.Context(), p)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve categories", err.Error())
		return
	}

	response.PaginatedJSON(c, http.StatusOK, "Categories retrieved successfully", categories, total, p.Page, p.Limit)
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	category, err := h.usecase.GetCategoryByPublicID(c.Request.Context(), publicID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch category", err.Error())
		return
	}

	if category == nil {
		response.Error(c, http.StatusNotFound, "Category not found", nil)
		return
	}

	response.JSON(c, http.StatusOK, "Category fetched successfully", category)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	category := &domain.Category{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Slug:        req.Slug,
		IconURL:     req.IconURL,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}

	if err := h.usecase.CreateCategory(c.Request.Context(), category); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Category created successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	category := &domain.Category{
		BaseEntity:  domain.BaseEntity{PublicID: publicID},
		ParentID:    req.ParentID,
		Name:        req.Name,
		Slug:        req.Slug,
		IconURL:     req.IconURL,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}

	if err := h.usecase.UpdateCategory(c.Request.Context(), category); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update category", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Category updated successfully", category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	if err := h.usecase.DeleteCategory(c.Request.Context(), publicID); err != nil {
		if err == domain.ErrEntityInUse {
			response.Error(c, http.StatusConflict, "Category is currently in use", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete category", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

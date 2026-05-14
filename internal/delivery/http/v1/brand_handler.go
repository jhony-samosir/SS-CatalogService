package v1

import (
	"net/http"
	"ss-catalog-service/internal/delivery/http/v1/dto"
	"ss-catalog-service/internal/domain"
	"ss-catalog-service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BrandHandler struct {
	usecase domain.BrandUsecase
}

func NewBrandHandler(u domain.BrandUsecase) *BrandHandler {
	return &BrandHandler{usecase: u}
}

func (h *BrandHandler) GetBrands(c *gin.Context) {
	var p domain.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p = domain.Pagination{Limit: 100, Offset: 0}
	}

	brands, err := h.usecase.GetBrands(c.Request.Context(), p)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch brands", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Brands fetched successfully", brands)
}

func (h *BrandHandler) GetBrand(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand ID", nil)
		return
	}

	brand, err := h.usecase.GetBrandByPublicID(c.Request.Context(), publicID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch brand", err.Error())
		return
	}

	if brand == nil {
		response.Error(c, http.StatusNotFound, "Brand not found", nil)
		return
	}

	response.JSON(c, http.StatusOK, "Brand fetched successfully", brand)
}

func (h *BrandHandler) CreateBrand(c *gin.Context) {
	var req dto.CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	brand := &domain.Brand{
		Name:        req.Name,
		Slug:        req.Slug,
		LogoURL:     req.LogoURL,
		WebsiteURL:  req.WebsiteURL,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := h.usecase.CreateBrand(c.Request.Context(), brand); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create brand", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Brand created successfully", brand)
}

func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand ID", nil)
		return
	}

	var req dto.UpdateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	brand := &domain.Brand{
		BaseEntity:  domain.BaseEntity{PublicID: publicID},
		Name:        req.Name,
		Slug:        req.Slug,
		LogoURL:     req.LogoURL,
		WebsiteURL:  req.WebsiteURL,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := h.usecase.UpdateBrand(c.Request.Context(), brand); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update brand", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Brand updated successfully", brand)
}

func (h *BrandHandler) DeleteBrand(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand ID", nil)
		return
	}

	if err := h.usecase.DeleteBrand(c.Request.Context(), publicID); err != nil {
		if err == domain.ErrEntityInUse {
			response.Error(c, http.StatusConflict, "Brand is currently in use", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete brand", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

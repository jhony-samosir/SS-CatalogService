package v1

import (
	"net/http"
	"ss-catalog-service/internal/delivery/http/v1/dto"
	"ss-catalog-service/internal/domain"
	"ss-catalog-service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WarehouseHandler struct {
	usecase domain.WarehouseUsecase
}

func NewWarehouseHandler(u domain.WarehouseUsecase) *WarehouseHandler {
	return &WarehouseHandler{usecase: u}
}

func (h *WarehouseHandler) GetWarehouses(c *gin.Context) {
	var p domain.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p = domain.Pagination{Limit: 100, Offset: 0}
	}

	warehouses, err := h.usecase.GetWarehouses(c.Request.Context(), p)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch warehouses", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Warehouses fetched successfully", warehouses)
}

func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var req dto.CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	warehouse := &domain.Warehouse{
		Name:        req.Name,
		Code:        req.Code,
		City:        req.City,
		Province:    req.Province,
		CountryCode: req.CountryCode,
		PostalCode:  req.PostalCode,
		Address:     req.Address,
		IsActive:    req.IsActive,
	}

	if err := h.usecase.CreateWarehouse(c.Request.Context(), warehouse); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create warehouse", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Warehouse created successfully", warehouse)
}

func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid warehouse ID", nil)
		return
	}

	var req dto.UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	warehouse := &domain.Warehouse{
		BaseEntity:  domain.BaseEntity{PublicID: publicID},
		Name:        req.Name,
		Code:        req.Code,
		City:        req.City,
		Province:    req.Province,
		CountryCode: req.CountryCode,
		PostalCode:  req.PostalCode,
		Address:     req.Address,
		IsActive:    req.IsActive,
	}

	if err := h.usecase.UpdateWarehouse(c.Request.Context(), warehouse); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update warehouse", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Warehouse updated successfully", warehouse)
}

func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid warehouse ID", nil)
		return
	}

	if err := h.usecase.DeleteWarehouse(c.Request.Context(), publicID); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete warehouse", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

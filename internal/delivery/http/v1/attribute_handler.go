package v1

import (
	"net/http"
	"ss-catalog-service/internal/delivery/http/v1/dto"
	"ss-catalog-service/internal/domain"
	"ss-catalog-service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttributeHandler struct {
	usecase domain.AttributeUsecase
}

func NewAttributeHandler(u domain.AttributeUsecase) *AttributeHandler {
	return &AttributeHandler{usecase: u}
}

func (h *AttributeHandler) GetAttributes(c *gin.Context) {
	var p domain.Pagination
	_ = c.ShouldBindQuery(&p)
	p.SetDefaults()

	attributes, total, err := h.usecase.GetAttributes(c.Request.Context(), p)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve attributes", err.Error())
		return
	}

	response.PaginatedJSON(c, http.StatusOK, "Attributes retrieved successfully", attributes, total, p.Page, p.Limit)
}

func (h *AttributeHandler) GetAttribute(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid attribute ID", nil)
		return
	}

	attribute, err := h.usecase.GetAttributeByPublicID(c.Request.Context(), publicID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch attribute", err.Error())
		return
	}

	if attribute == nil {
		response.Error(c, http.StatusNotFound, "Attribute not found", nil)
		return
	}

	response.JSON(c, http.StatusOK, "Attribute fetched successfully", attribute)
}

func (h *AttributeHandler) CreateAttribute(c *gin.Context) {
	var req dto.CreateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	attribute := &domain.ProductAttribute{
		Name:      req.Name,
		Code:      req.Code,
		InputType: domain.AttributeInputType(req.InputType),
		IsVariant: req.IsVariant,
		SortOrder: req.SortOrder,
	}

	if err := h.usecase.CreateAttribute(c.Request.Context(), attribute); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create attribute", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Attribute created successfully", attribute)
}

func (h *AttributeHandler) UpdateAttribute(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid attribute ID", nil)
		return
	}

	var req dto.UpdateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	attribute := &domain.ProductAttribute{
		BaseEntity: domain.BaseEntity{PublicID: publicID},
		Name:       req.Name,
		Code:       req.Code,
		InputType:  domain.AttributeInputType(req.InputType),
		IsVariant:  req.IsVariant,
		SortOrder:  req.SortOrder,
	}

	if err := h.usecase.UpdateAttribute(c.Request.Context(), attribute); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update attribute", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Attribute updated successfully", attribute)
}

func (h *AttributeHandler) DeleteAttribute(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid attribute ID", nil)
		return
	}

	if err := h.usecase.DeleteAttribute(c.Request.Context(), publicID); err != nil {
		if err == domain.ErrEntityInUse {
			response.Error(c, http.StatusConflict, "Attribute is currently in use", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete attribute", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

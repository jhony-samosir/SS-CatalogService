package v1

import (
	"net/http"
	"ss-catalog-service/internal/domain"
	"ss-catalog-service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TagHandler struct {
	usecase domain.TagUsecase
}

func NewTagHandler(u domain.TagUsecase) *TagHandler {
	return &TagHandler{usecase: u}
}

func (h *TagHandler) GetTags(c *gin.Context) {
	var p domain.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		p.SetDefaults()
	}
	p.SetDefaults()

	tags, total, err := h.usecase.GetTags(c.Request.Context(), p)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch tags", err.Error())
		return
	}

	response.PaginatedJSON(c, http.StatusOK, "Tags fetched successfully", tags, total, p.Page, p.Limit)
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var tag domain.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := h.usecase.CreateTag(c.Request.Context(), &tag); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create tag", err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Tag created successfully", tag)
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid tag ID", nil)
		return
	}

	var tag domain.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		response.ValidationError(c, err)
		return
	}

	tag.PublicID = publicID
	if err := h.usecase.UpdateTag(c.Request.Context(), &tag); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update tag", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Tag updated successfully", tag)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	publicID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid tag ID", nil)
		return
	}

	if err := h.usecase.DeleteTag(c.Request.Context(), publicID); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete tag", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

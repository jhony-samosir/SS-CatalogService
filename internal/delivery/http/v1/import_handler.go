package v1

import (
	"net/http"
	"strconv"
	"ss-catalog-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImportHandler struct {
	usecase domain.ImportUsecase
}

func NewImportHandler(u domain.ImportUsecase) *ImportHandler {
	return &ImportHandler{usecase: u}
}

type TriggerImportRequest struct {
	FileURL string `json:"file_url" binding:"required,url"`
	JobType string `json:"job_type" binding:"required"`
}

func (h *ImportHandler) TriggerImport(c *gin.Context) {
	var req TriggerImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	job, err := h.usecase.TriggerImport(c.Request.Context(), req.FileURL, req.JobType, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to trigger import"})
		return
	}

	c.JSON(http.StatusAccepted, job)
}

func (h *ImportHandler) GetJobStatus(c *gin.Context) {
	publicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	job, err := h.usecase.GetJobStatus(c.Request.Context(), publicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *ImportHandler) GetImportJobs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	jobs, total, err := h.usecase.GetAllJobs(c.Request.Context(), domain.Pagination{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get import jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"items":       jobs,
			"total_count": total,
		},
	})
}

package v1

import (
	"net/http"
	"strconv"
	"ss-catalog-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type DigitalHandler struct {
	usecase domain.DigitalUsecase
}

func NewDigitalHandler(u domain.DigitalUsecase) *DigitalHandler {
	return &DigitalHandler{usecase: u}
}

type AddLicensesRequest struct {
	ProductID int      `json:"product_id" binding:"required"`
	Keys      []string `json:"keys" binding:"required,min=1"`
}

func (h *DigitalHandler) AddLicenses(c *gin.Context) {
	var req AddLicensesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.usecase.AddLicenses(c.Request.Context(), req.ProductID, req.Keys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add licenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "licenses added successfully"})
}

func (h *DigitalHandler) GetDigitalDetails(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	files, licenseCount, err := h.usecase.GetDigitalDetails(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get digital details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files":         files,
		"license_count": licenseCount,
	})
}

type UploadFileRequest struct {
	ProductID     int    `json:"product_id" binding:"required"`
	FileName      string `json:"file_name" binding:"required"`
	FilePath      string `json:"file_path" binding:"required"`
	FileSizeBytes int64  `json:"file_size_bytes"`
	MimeType      string `json:"mime_type"`
	Version       string `json:"version"`
}

func (h *DigitalHandler) UploadFile(c *gin.Context) {
	var req UploadFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.usecase.UploadDigitalProduct(c.Request.Context(), &domain.DigitalFile{
		ProductID:     req.ProductID,
		FileName:      req.FileName,
		FilePath:      req.FilePath,
		FileSizeBytes: req.FileSizeBytes,
		MimeType:      req.MimeType,
		Version:       req.Version,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload digital product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "digital product uploaded successfully"})
}

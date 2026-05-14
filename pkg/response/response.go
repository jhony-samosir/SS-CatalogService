package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type PaginatedData struct {
	Items      interface{} `json:"items"`
	TotalCount int64       `json:"total_count"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
}

func JSON(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func PaginatedJSON(c *gin.Context, status int, message string, items interface{}, total int64, page, limit int) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data: PaginatedData{
			Items:      items,
			TotalCount: total,
			Page:       page,
			Limit:      limit,
		},
	})
}

func Error(c *gin.Context, status int, message string, errs interface{}) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Errors:  errs,
	})
}

func ValidationError(c *gin.Context, err error) {
	errs := make(map[string]string)
	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, f := range verrs {
			errs[f.Field()] = f.Tag()
		}
	} else {
		errs["error"] = err.Error()
	}

	Error(c, http.StatusBadRequest, "Validation failed", errs)
}

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

type PaginationResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

func Paginated(c *gin.Context, data interface{}, total int64, page, limit int) {
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	c.JSON(http.StatusOK, PaginationResponse{
		Success: true,
		Data:    data,
		Meta: PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	})
}

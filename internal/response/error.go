package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
	})
}

func ValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Message: "Validasi gagal",
		Errors:  errors,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Tidak terautentikasi"
	}
	Error(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Akses ditolak"
	}
	Error(c, http.StatusForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Data tidak ditemukan"
	}
	Error(c, http.StatusNotFound, message)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message)
}

func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "Terjadi kesalahan pada server"
	}
	Error(c, http.StatusInternalServerError, message)
}

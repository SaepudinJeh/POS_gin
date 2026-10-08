package middlewares

import (
	"log"

	"POS/internal/response"

	"github.com/gin-gonic/gin"
)

// Recovery menangkap panic dan mengembalikan response 500 yang seragam.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v | %s %s", err, c.Request.Method, c.Request.URL.Path)
				response.InternalError(c, "Terjadi kesalahan pada server")
				c.Abort()
			}
		}()

		c.Next()
	}
}

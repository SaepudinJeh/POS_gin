package middlewares

import (
	"POS/internal/response"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "Role tidak ditemukan di token")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			response.Forbidden(c, "Format role tidak valid")
			c.Abort()
			return
		}

		for _, allowed := range roles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Akses ditolak")
		c.Abort()
	}
}

package routes

import (
	"net/http"

	"POS/internal/handlers"
	"POS/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth     *handlers.AuthHandler
	Category *handlers.CategoryHandler
}

func SetupRouter(h *Handlers) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	// ---------- Public routes ----------
	RegisterAuthRoutes(api, h.Auth)

	// ---------- Protected routes ----------
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware())

	// Profile (sementara inline, nanti bisa dipindah ke UserHandler)
	protected.GET("/profile", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"user_id": userID,
				"role":    role,
			},
		})
	})

	RegisterCategoryRoutes(protected, h.Category)
	// nanti: RegisterProductRoutes(protected, h.Product)

	return r
}

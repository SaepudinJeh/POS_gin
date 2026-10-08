package routes

import (
	"POS/internal/handlers"
	"POS/internal/middlewares"
	"POS/internal/response"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth     *handlers.AuthHandler
	Category *handlers.CategoryHandler
	Product  *handlers.ProductHandler
}

func SetupRouter(h *Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middlewares.Recovery())

	api := r.Group("/api/v1")

	// Public routes
	RegisterAuthRoutes(api, h.Auth)

	// Protected routes
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware())

	protected.GET("/profile", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")
		response.Success(c, gin.H{
			"user_id": userID,
			"role":    role,
		})
	})

	RegisterCategoryRoutes(protected, h.Category)
	RegisterProductRoutes(protected, h.Product)

	return r
}

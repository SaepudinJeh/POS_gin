package routes

import (
	"POS/internal/handlers"
	"POS/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(rg *gin.RouterGroup, h *handlers.ProductHandler) {
	products := rg.Group("/products")
	{
		products.GET("", h.GetAll)
		products.GET("/:id", h.GetByID)

		admin := products.Group("")
		admin.Use(middlewares.RequireRole("admin"))
		{
			admin.POST("", h.Create)
			admin.PUT("/:id", h.Update)
			admin.DELETE("/:id", h.Delete)
		}
	}
}

package routes

import (
	"POS/internal/handlers"
	"POS/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(rg *gin.RouterGroup, h *handlers.CategoryHandler) {
	categories := rg.Group("/categories")
	{
		// Semua user yang login bisa lihat
		categories.GET("", h.GetAll)
		categories.GET("/:id", h.GetByID)

		// Hanya admin
		admin := categories.Group("")
		admin.Use(middlewares.RequireRole("admin"))
		{
			admin.POST("", h.Create)
			admin.PUT("/:id", h.Update)
			admin.DELETE("/:id", h.Delete)
		}
	}
}

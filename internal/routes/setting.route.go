package routes

import (
	"POS/internal/handlers"
	"POS/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterSettingRoutes(rg *gin.RouterGroup, h *handlers.SettingHandler) {
	settings := rg.Group("/settings")
	settings.Use(middlewares.RequireRole("admin")) // admin only
	{
		settings.GET("", h.GetAll)
		settings.PUT("", h.BulkUpdate)
		settings.PUT("/:key", h.Update)
	}
}

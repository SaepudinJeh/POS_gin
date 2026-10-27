package routes

import (
	"POS/internal/handlers"
	"POS/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterReportRoutes(rg *gin.RouterGroup, h *handlers.ReportHandler) {
	reports := rg.Group("/reports")
	reports.Use(middlewares.RequireRole("admin")) // admin only
	{
		reports.GET("/summary", h.Summary)
		reports.GET("/daily", h.Daily)
		reports.GET("/by-cashier", h.ByCashier)
		reports.GET("/by-product", h.ByProduct)
		reports.GET("/by-category", h.ByCategory)
	}
}

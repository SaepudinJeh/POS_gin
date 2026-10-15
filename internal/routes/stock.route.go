package routes

import (
	"POS/internal/handlers"
	"POS/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterStockRoutes(rg *gin.RouterGroup, h *handlers.StockHandler) {
	stock := rg.Group("/stock")
	{
		stock.GET("/mutations", h.GetAll)
		stock.GET("/mutations/:id", h.GetByID)

		// Hanya admin yang boleh catat mutasi
		admin := stock.Group("")
		admin.Use(middlewares.RequireRole("admin"))
		{
			admin.POST("/mutations", h.Create)
		}
	}

	// Riwayat mutasi per produk
	products := rg.Group("/products")
	{
		products.GET("/:id/mutations", h.GetByProductID)
	}
}

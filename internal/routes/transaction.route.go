package routes

import (
	"POS/internal/handlers"
	"POS/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterTransactionRoutes(rg *gin.RouterGroup, h *handlers.TransactionHandler) {
	txns := rg.Group("/transactions")
	{
		txns.GET("", h.GetAll)
		txns.GET("/:id", h.GetByID)
		txns.POST("", h.Create)
		txns.POST("/:id/pay", h.Pay)

		// Void: hanya admin
		admin := txns.Group("")
		admin.Use(middlewares.RequireRole("admin"))
		{
			admin.POST("/:id/void", h.Void)
		}
	}
}

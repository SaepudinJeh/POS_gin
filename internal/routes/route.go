package routes

import (
	"POS/internal/handlers"
	"POS/internal/middlewares"
	"POS/internal/response"

	"github.com/gin-gonic/gin"

	_ "POS/docs" // <-- akan dibuat di STEP 7

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	Auth        *handlers.AuthHandler
	Category    *handlers.CategoryHandler
	Product     *handlers.ProductHandler
	Stock       *handlers.StockHandler
	Transaction *handlers.TransactionHandler
	Setting     *handlers.SettingHandler
	Report      *handlers.ReportHandler
}

func SetupRouter(h *Handlers) *gin.Engine {
	r := gin.New()

	// 1. Pasang semua middleware global
	middlewares.RegisterGlobal(r)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 2. Config engine
	r.HandleMethodNotAllowed = true
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})

	// 3. Route
	api := r.Group("/api/v1")
	RegisterAuthRoutes(api, h.Auth)

	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware())

	RegisterCategoryRoutes(protected, h.Category)
	RegisterProductRoutes(protected, h.Product)
	RegisterStockRoutes(protected, h.Stock)
	RegisterTransactionRoutes(protected, h.Transaction)
	RegisterSettingRoutes(protected, h.Setting)
	RegisterReportRoutes(protected, h.Report)

	// 4. Fallback
	r.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "Route tidak ditemukan")
	})
	r.NoMethod(func(c *gin.Context) {
		response.Error(c, 405, "Method tidak diizinkan")
	})

	return r
}

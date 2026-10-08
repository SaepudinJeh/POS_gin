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

	// 1. Pasang semua middleware global
	middlewares.RegisterGlobal(r)

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

	// 4. Fallback
	r.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "Route tidak ditemukan")
	})
	r.NoMethod(func(c *gin.Context) {
		response.Error(c, 405, "Method tidak diizinkan")
	})

	return r
}

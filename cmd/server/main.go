package main

import (
	"POS/internal/configs"
	"POS/internal/handlers"
	"POS/internal/middlewares"
	"POS/internal/repositories"
	"POS/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, pakai environment system")
	}

	configs.ConnectDatabase()

	userRepo := repositories.NewUserRepository(configs.DB)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()

	authRoutes := r.Group("/api/v1/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	protected := r.Group("/api/v1")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{"user_id": userID})
		})
	}

	r.TrustedPlatform = gin.PlatformFlyIO
	r.TrustedPlatform = gin.PlatformCloudflare
	r.TrustedPlatform = gin.PlatformGoogleAppEngine

	r.HandleMethodNotAllowed = true

	r.Use(middlewares.RateLimiter())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middlewares.ErrorHandler())

	// handler global error
	r.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Route not found",
		})
	})

	r.NoMethod(func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{
			"success": false,
			"message": "Method not allowed",
		})
	})

	r.Run()
}

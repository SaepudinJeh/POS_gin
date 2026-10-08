package main

import (
	"POS/internal/configs"
	"POS/internal/container"
	"POS/internal/middlewares"
	"POS/internal/routes"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, pakai environment system")
	}

	configs.ConnectDatabase()

	r := gin.Default()

	c := container.NewContainer(configs.DB)

	r = routes.SetupRouter(&routes.Handlers{
		Auth:     c.Auth,
		Category: c.Category,
	})

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notspace", func(fl validator.FieldLevel) bool {
			return !strings.Contains(fl.Field().String(), " ")
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

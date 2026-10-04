package main

import (
	"POS/internal/middlewares"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.TrustedPlatform = gin.PlatformFlyIO
	r.TrustedPlatform = gin.PlatformCloudflare
	r.TrustedPlatform = gin.PlatformGoogleAppEngine

	r.HandleMethodNotAllowed = true

	r.Use(middlewares.RateLimiter())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://example.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middlewares.ErrorHandler())

	r.GET("/ping", func(ctx *gin.Context) {
		fmt.Printf("ClientIP: %s\n", ctx.ClientIP())
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Pong",
		})
	})

	// handler global error
	r.NoRoute(func(ctx *gin.Context) {
		fmt.Println(">>> NoRoute triggered:", ctx.Request.URL.Path) // debug
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

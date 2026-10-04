package main

import (
	"POS/internal/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.TrustedPlatform = gin.PlatformFlyIO

	r.HandleMethodNotAllowed = true

	r.Use(middlewares.ErrorHandler())

	r.GET("/pingg", func(ctx *gin.Context) {
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

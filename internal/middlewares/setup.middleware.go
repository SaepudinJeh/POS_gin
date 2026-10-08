package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RegisterGlobal memasang semua middleware global ke engine.
// Dipanggil SEKALI di SetupRouter, sebelum route didaftarkan.
func RegisterGlobal(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(Recovery())
	r.Use(RateLimiter())
	r.Use(CORS())
	r.Use(ErrorHandler())
}

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
}

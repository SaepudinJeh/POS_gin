package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thienel/tlog"
)

func RegisterGlobal(r *gin.Engine) {
	// tlog middleware untuk request logging
	r.Use(tlog.GinMiddleware(
		tlog.WithSkipPaths("/health", "/metrics", "/favicon.ico"),
		tlog.WithMaskPatterns(
			`(?i)password`,
			`(?i)token`,
			`(?i)secret`,
			`(?i)authorization`,
		),
	))

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

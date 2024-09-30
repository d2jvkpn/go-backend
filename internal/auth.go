package internal

import (
	// "fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(origins []string, maxAges ...time.Duration) gin.HandlerFunc {
	maxAge := 12 * time.Hour
	if len(maxAges) > 0 {
		maxAge = maxAges[0]
	}

	return cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "OPTIONS", "HEAD"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"x-client",
		},
		ExposeHeaders: []string{
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Content-Disposition",
		},
		AllowWildcard:    true,
		AllowCredentials: true,
		// AllowOriginFunc:  func(origin string) bool { return origin == "https://github.com" },
		MaxAge: maxAge,
	})
}

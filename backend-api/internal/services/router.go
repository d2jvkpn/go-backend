package services

import (
	// "fmt"

	"backend-api/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadOpen(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	group := router.Group("/api/v1/open", handlers...)

	group.POST("/account/login", accountLogin)
}

func LoadAuth(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	group := router.Group("/api/v1/auth", handlers...)

	group.GET("/hello", func(ctx *gin.Context) {
		middlewares.JsonOK(ctx)
	})
}

func LoadWebsocket(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	ws := router.Group("/socket", handlers...)
	ws.GET("/talk", talk)
}

package services

import (
	// "fmt"

	"github.com/gin-gonic/gin"
)

func LoadOpen(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	_ = router.Group("/api/v1/open", handlers...)
}

func LoadWebsocket(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	ws := router.Group("/socket", handlers...)
	ws.GET("/talk", talk)
}

package services

import (
	// "fmt"

	// "backend-api/pkg/structs"

	"github.com/gin-gonic/gin"
)

func LoadOpen(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	group := router.Group("/api/v1/open", handlers...)

	group.POST("/account/login", accountLogin)
}

/*
	group.GET("/hello", func(ctx *gin.Context) {
		structs.JsonOK(ctx)
	})
*/

func LoadAuth(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	group := router.Group("/api/v1/auth", handlers...)

	group.POST("/account/logout", accountLogout)
	group.POST("/account/change_password", accountChangePassword)
	group.GET("/account/query_accounts", queryAccounts)
}

func LoadWebsocket(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	ws := router.Group("/api/v1/socket", handlers...)
	ws.GET("/talk", talk)
}

package api

import (
	// "fmt"

	"github.com/gin-gonic/gin"
)

func LoadOpen(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	_ = router.Group("/api/v1/open", handlers...)
}

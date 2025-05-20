package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetAPI(ctx *gin.Context) string {
	return fmt.Sprintf("%s@%s", ctx.Request.Method, ctx.Request.URL.Path)
}

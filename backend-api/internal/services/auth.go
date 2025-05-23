package services

import (
	"fmt"
	"slices"

	"backend-api/pkg/structs"

	"github.com/gin-gonic/gin"
)

func AllowLevels(levels []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		level := ctx.GetString("level")
		if !slices.Contains(levels, level) {
			err := structs.NotPermited(fmt.Errorf("no permited for %s", level))
			structs.JsonErr(ctx, err)
			ctx.Abort()
		}

		ctx.Next()
	}
}

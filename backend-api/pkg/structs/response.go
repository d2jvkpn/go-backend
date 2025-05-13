package structs

import (
	// "fmt"
	"net/http"

	"backend-api/pkg/erri"

	"github.com/d2jvkpn/errx"
	"github.com/gin-gonic/gin"
)

func JSONErr(ctx *gin.Context, err *errx.ErrX) {
	status := erri.ErrStatus(err)
	ctx.Set("Error", err)

	ctx.JSON(status, gin.H{
		"requestId": ctx.GetString("RequestId"),
		"kind":      err.Kind,
		"code":      err.Code,
		"msg":       err.Msg,
	})

	return
}

func JSONOK(ctx *gin.Context, data ...any) {
	requestId := ctx.GetString("RequestId")

	if len(data) == 0 {
		data = []any{gin.H{}}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"requestId": requestId,
		"code":      "ok",
		"data":      data[0],
	})
}

package middlewares

import (
	"bytes"
	"net/http"

	"github.com/d2jvkpn/errx"
	"github.com/gin-gonic/gin"
)

func ErrStatus(err *errx.ErrX) (status int) {
	switch err.Code {
	case "BadRequest":
		status = http.StatusBadRequest
	case "Unauthorized":
		status = http.StatusUnauthorized
	case "NoPermission":
		status = http.StatusForbidden
	case "NoRoute":
		status = http.StatusNotFound
	case "NoAcceptable":
		status = http.StatusNotAcceptable
	case "BizError":
		status = http.StatusConflict
	case "InternalError":
		status = http.StatusInternalServerError
	case "ServiceUnavailable":
		status = http.StatusServiceUnavailable
	default:
		status = 599
	}

	return status
}

type ResponseOK struct {
	// string: uuid
	RequestId string `json:"requestId" example:"3cc643bd-7f85-493e-8324-6e491db7b3d8" extensions:"x-order=01"`

	// string: OK
	Code string `json:"code" example:"OK" extensions:"x-order=02"`

	// any: response data
	Data any `json:"data" swaggertype:"object,string" example:"answer:hello,value:42" extensions:"x-order=03"`
}

type ResponseItem[T any] struct {
	// string: uuid
	RequestId string `json:"requestId" example:"3cc643bd-7f85-493e-8324-6e491db7b3d8" extensions:"x-order=01"`

	// string: OK
	Code string `json:"code" example:"OK" extensions:"x-order=02"`

	Data struct {
		Item T `json:"item" swaggertype:"object,string" example:"answer:hello,value:42"`
	} `json:"data,omitempty" extensions:"x-order=03"`
}

type ResponsePage[T any] struct {
	// string: uuid
	RequestId string `json:"requestId" example:"3cc643bd-7f85-493e-8324-6e491db7b3d8" extensions:"x-order=01"`

	// string: OK
	Code string `json:"code" example:"OK" extensions:"x-order=02"`

	Data struct {
		PageIndex uint `json:"pageIndex" example:"1"`
		PageSize  uint `json:"pageSize" example:"30"`
		Total     uint `json:"total" example:"100"`
		Items     []T  `json:"items"`
	} `json:"data" extensions:"x-order=03"`
}

type ResponseErr struct {
	// string: uuid
	RequestId string `json:"requestId" example:"3cc643bd-7f85-493e-8324-6e491db7b3d8" extensions:"x-order=01"`

	// string: error code(BadRequest....)
	Code string `json:"code" example:"BadRequest" extensions:"x-order=02"`

	// string: error kind(invalid_parameter)
	Kind string `json:"kind" example:"invalid_parameter" extensions:"x-order=03"`

	// Option<string>: notification msg
	Msg string `json:"msg" example:"no account id" extensions:"x-order=04"`
}

func JsonErr(ctx *gin.Context, err *errx.ErrX) {
	status := ErrStatus(err)
	ctx.Set("error", err)

	ctx.JSON(status, ResponseErr{
		RequestId: ctx.GetString("requestId"),
		Code:      err.Code,
		Kind:      err.Kind,
		Msg:       err.Msg,
	})

	return
}

func JsonOK(ctx *gin.Context, data ...any) {
	requestId := ctx.GetString("requestId")

	if len(data) == 0 {
		data = []any{gin.H{}}
	}
	ctx.JSON(http.StatusOK, ResponseOK{RequestId: requestId, Code: "OK", Data: data[0]})
}

func ResponseFile(ctx *gin.Context, buf *bytes.Buffer, filename, typ string) {
	var contextType string

	switch typ {
	case "xls", "xlsx":
		contextType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "doc", "docx":
		contextType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "pdf":
		contextType = "application/pdf"
	default:
		contextType = "application/octet-stream"
	}

	// ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)

	ctx.Data(http.StatusOK, contextType, buf.Bytes())
}

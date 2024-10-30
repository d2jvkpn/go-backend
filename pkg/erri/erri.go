package erri

import (
	// "fmt"

	"github.com/d2jvkpn/errx"
)

func Validation(e error, msg string, args ...any) *errx.ErrX {
	return errx.New(e).
		WithKind("invalid_parameter").
		WithCode("validation").
		WithMsg(msg, args...).
		WithCaller(2)
}

func InternalErr(e error, kind string) *errx.ErrX {
	return errx.New(e).
		WithKind("internal_error").
		WithCode(kind).
		WithMsg("an internal error occured").
		WithCaller(2)
}

func BizErr(e error, kind string) *errx.ErrX {
	return errx.New(e).
		WithKind("biz_error").
		WithCode(kind).
		WithCaller(2)
}

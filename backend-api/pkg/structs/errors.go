// default: ErrX.Kind == ErrX.Code
package structs

import (
	// "fmt"

	"github.com/d2jvkpn/errx"
)

// no route
func NoRoute() *errx.ErrX {
	return errx.Eee().WithKind("no_route").WithCode("no_route")
}

// invalid
// func Invalid(e error, msg string, args ...any) *errx.ErrX {
func Invalid(e error) *errx.ErrX {
	return errx.New(e).WithKind("invalid").WithCode("invalid").WithCaller(2)
}

// incorrect
func Incorrect(e error) *errx.ErrX {
	return errx.New(e).WithKind("incorrect").WithCode("incorrect").WithCaller(2)
}

// bind error
func BindError(e error) *errx.ErrX {
	return errx.New(e).WithKind("bind_error").WithCode("bind_error").WithCaller(2)
}

// authorization error
func AuthError(e error) *errx.ErrX {
	return errx.New(e).WithKind("auth_err").WithCode("auth_err")
}

// not permited
func NotPermited(e error) *errx.ErrX {
	return errx.New(e).WithKind("not_permited").WithCode("not_permited")
}

// biz error
func bizError(e error) *errx.ErrX {
	return errx.New(e).WithKind("biz_error").WithMsg("an biz error occured")
}

// biz error
func BizError(e error) *errx.ErrX {
	return bizError(e).WithCode("biz_error").WithCaller(2)
}

// internal error
func internalError(e error) *errx.ErrX {
	return errx.New(e).WithKind("internal_error").WithMsg("an internal error occured")
}

// internal error
func InternalError(e error) *errx.ErrX {
	return internalError(e).WithCode("internal_error").WithCaller(2)
}

// unavailable
func Unavailable(e error) *errx.ErrX {
	return errx.New(e).
		WithKind("unavailable").
		WithCode("unavailable").
		WithMsg("an service is unavailable").
		WithCaller(2)
}

// unknown error
func UnknownError(e error) *errx.ErrX {
	return errx.New(e).
		WithKind("unknown_error").
		WithCode("unknown_error").
		WithMsg("an unknown error occured").
		WithCaller(2)
}

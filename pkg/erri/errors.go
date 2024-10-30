package erri

import (
	// "fmt"

	"github.com/d2jvkpn/errx"
)

// no route
func NoRoute() *errx.ErrX {
	return errx.Eee().
		WithKind("no_route").
		WithCode("no_route")
}

// invalid
// func Invalid(e error, msg string, args ...any) *errx.ErrX {
func Invalid(e error) *errx.ErrX {
	return errx.New(e).
		WithKind("invalid").
		WithCode("invalid").
		WithCaller(2)
}

// incorrect
func Incorrect(e error) *errx.ErrX {
	return errx.New(e).
		WithKind("incorrect").
		WithCode("incorrect").
		WithCaller(2)
}

// bind error
func BindErr(e error, code string) *errx.ErrX {
	return errx.New(e).
		WithKind("bind_error").
		WithCode(code).
		WithCaller(2)
}

// authorization error
func AuthErr(e error, code string) *errx.ErrX {
	return errx.New(e).
		WithKind("authorization_error").
		WithCode(code)
}

// not permited
func NotPermited(e error, code string) *errx.ErrX {
	return errx.New(e).
		WithKind("not_permited").
		WithCode(code)
}

// biz error
func BizErr(e error, kind string) *errx.ErrX {
	return errx.New(e).
		WithKind("biz_error").
		WithCode(kind).
		WithCaller(2)
}

// internal error
func InternalErr(e error, kind string) *errx.ErrX {
	return errx.New(e).
		WithKind("internal_error").
		WithCode(kind).
		WithMsg("an internal error occured").
		WithCaller(2)
}

// unavailable
func Unavailable(e error, kind string) *errx.ErrX {
	return errx.New(e).
		WithKind("unavailable").
		WithCode(kind).
		WithMsg("an service is unavailable").
		WithCaller(2)
}

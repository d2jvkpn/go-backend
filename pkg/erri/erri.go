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

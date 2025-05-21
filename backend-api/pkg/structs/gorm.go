package structs

import (
	// "fmt"

	"backend-api/pkg/infra"

	"github.com/d2jvkpn/errx"
)

func PgNotFound(e error) (err *errx.ErrX) {
	if infra.PgNotFound(e) {
		return bizError(e).WithCode("not_found").WithCaller(2)
	}

	return internalError(e).WithCode("database").WithCaller(2)
}

func PgIsUnique(e error) (err *errx.ErrX) {
	if infra.PgUniqueViolation(e) {
		return bizError(e).WithCode("already_exists").WithCaller(2)
	}

	return internalError(e).WithCode("database").WithCaller(2)
}

func PgIsNotFoundOrUnique(e error) (err *errx.ErrX) {
	if infra.PgNotFound(e) {
		return bizError(e).WithCode("not_found").WithCaller(2)
	}

	if infra.PgUniqueViolation(e) {
		return bizError(e).WithCode("already_exists").WithCaller(2)
	}

	return internalError(e).WithCode("database").WithCaller(2)
}

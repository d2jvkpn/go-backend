package erri

import (
	// "fmt"
	"net/http"

	"github.com/d2jvkpn/errx"
)

func ErrStatus(err *errx.ErrX) (status int) {
	switch err.Kind {
	case "no_route": // 404
		status = http.StatusNotFound
	case "invalid": // 400
		status = http.StatusBadRequest
	case "incorrect": // 400
		status = http.StatusBadRequest
	case "bind_error": // 400
		status = http.StatusBadRequest
	case "authorization_error": // 401
		status = http.StatusUnauthorized
	case "not_permited": // 403
		status = http.StatusForbidden
	case "biz_error": // 409
		status = http.StatusConflict
	case "internal_error": // 500
		status = http.StatusInternalServerError
	case "unavailable": // 503
		status = http.StatusServiceUnavailable
	default:
		status = 599
	}

	return status
}

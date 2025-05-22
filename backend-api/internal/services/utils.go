package services

import (
	"fmt"
	"strconv"
	"time"

	"backend-api/pkg/structs"
	"backend-api/pkg/utils"

	"github.com/d2jvkpn/errx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	_UUID_Null uuid.UUID
)

/*
func GetAuthAccount(ctx *gin.Context) (auth *structs.AuthAccount, err *errx.ErrX) {
	var e error

	//fmt.Printf("\n??? accountId=%s, tokenId=%s\n", ctx.GetString("accountId"), ctx.GetString("tokenId"))

	auth, e = structs.NewAuthAccount(
		ctx.GetString("accountId"),
		ctx.GetString("level"),
		ctx.GetString("tokenId"),
		ctx.GetString("platform"),
	)

	if e != nil {
		return nil, structs.InternalError(e).WithCode("context_no_value")
	}

	return auth, nil
}
*/

func GetAccountId(ctx *gin.Context) (accountId uuid.UUID, err *errx.ErrX) {
	var e error

	if accountId, e = utils.UUIDFromString(ctx.GetString("accountId")); e != nil {
		return accountId, structs.InternalError(e).WithCode("context_no_value")
	}

	return accountId, nil
}

// parameter must exists and not null
func QueryUUID(ctx *gin.Context, key string) (id uuid.UUID, err *errx.ErrX) {
	var (
		ok    bool
		value string
		e     error
	)

	if value, ok = ctx.GetQuery(key); !ok {
		return id, structs.Invalid(fmt.Errorf("no parameter: %s", key)).WithCode("no_parameter")
	}

	if id, e = utils.UUIDFromString(value); e != nil {
		return id, structs.Invalid(e).WithCode("parse_failed")
	}

	if utils.UUIDIsNull(id) {
		return id, structs.Invalid(utils.ErrIdIsNull)
	}

	return id, nil
}

// parameter must exists and not null
func QueryUUIDs(ctx *gin.Context, key string) (ids []uuid.UUID, err *errx.ErrX) {
	var (
		ok     bool
		values []string
		e      error
	)

	if values, ok = ctx.GetQueryArray(key); !ok {
		return nil, structs.Invalid(fmt.Errorf("no parameter: %s", key)).WithCode("no_parameter")
	}

	if ids, e = utils.UUIDFromStrings(values); e != nil {
		return nil, structs.Invalid(e).WithCode("parse_failed")
	}

	for i := range ids {
		if utils.UUIDIsNull(ids[i]) {
		}
	}

	return ids, nil
}

// parameter must exists
func QueryBool(ctx *gin.Context, key string) (ans bool, err *errx.ErrX) {
	var (
		ok    bool
		value string
		e     error
	)

	if value, ok = ctx.GetQuery(key); !ok {
		return false, structs.Invalid(fmt.Errorf("no parameter: %s", key)).WithCode("no_parameter")
	}

	if ans, e = strconv.ParseBool(value); e != nil {
		return false, structs.Invalid(e).WithCode("parse_failed")
	}

	return ans, nil
}

// parameter must exists
func QueryString(ctx *gin.Context, key string) (str string, err *errx.ErrX) {
	var ok bool

	if str, ok = ctx.GetQuery(key); !ok {
		return "", structs.Invalid(fmt.Errorf("no parameter: %s", key)).WithCode("no_parameter")
	}

	return str, nil
}

// parameter must exists
func QueryStrings(ctx *gin.Context, key string) (strs []string, err *errx.ErrX) {
	var ok bool

	if strs, ok = ctx.GetQueryArray(key); !ok {
		return nil, structs.Invalid(fmt.Errorf("no parameter: %s", key)).WithCode("no_parameter")
	}

	// TODO: ?? empty strings

	return strs, nil
}

// get date from url
func QueryDate(ctx *gin.Context, key string) (date string, err *errx.ErrX) {
	var (
		ok bool
		e  error
	)

	if date, ok = ctx.GetQuery(key); !ok {
		return "", structs.Invalid(fmt.Errorf("no parameter: %s", key)).WithCode("no_parameter")
	}

	if _, e = time.ParseInLocation(time.DateOnly, date, time.Local); e != nil {
		return "", structs.Invalid(e).WithCode("parse_failed")
	}

	return date, nil
}

func BindJSON[T any](ctx *gin.Context, value *T) (err *errx.ErrX) {
	var e error

	if e = ctx.BindJSON(value); e != nil {
		return structs.BindError(e).WithCode("bind_json")
	}

	return nil
}

func BindQuery[T any](ctx *gin.Context, value *T) (err *errx.ErrX) {
	var e error

	if e = ctx.BindQuery(value); e != nil {
		return structs.BindError(e).WithCode("bind_query")
	}

	return nil
}

func BindQueryJSON[T any](ctx *gin.Context, value *T) (err *errx.ErrX) {
	var e error

	if e = ctx.BindQuery(value); e != nil {
		return structs.BindError(e).WithCode("bind_query")
	}

	if e = ctx.BindJSON(value); e != nil {
		return structs.BindError(e).WithCode("bind_json")
	}

	return nil
}

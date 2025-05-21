package internal

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"backend-api/internal/settings"
	"backend-api/pkg/structs"

	"github.com/d2jvkpn/errx"
	"github.com/d2jvkpn/gotk/ginx"
	"github.com/gin-gonic/gin"
)

type HandleJwt func(context.Context, *ginx.JwtData) *errx.ErrX

func AllowLevels(levels ...string) HandleJwt {
	e := fmt.Errorf("allowed levels: %v", levels)

	return func(ctx context.Context, data *ginx.JwtData) (err *errx.ErrX) {
		if len(levels) > 0 && !slices.Contains(levels, data.Data["level"]) {
			err = structs.NotPermited(e).WithCode("no_allowed_level").WithMsg("no allowed level")
			return err
		}

		return nil
	}
}

func CacheUpdateToken(ctx context.Context, d *ginx.JwtData) (err *errx.ErrX) {
	// check if cache token enabled or not internal
	err = settings.CacheUpdateToken(ctx, fmt.Sprintf("login:%s:%s", d.Data["platform"], d.Subject), d.ID)

	if err != nil {
		// fmt.Printf("==> CacheUpdateToken: %v\n", err)
		return err
	}

	return nil
}

func Auth(funcs ...HandleJwt) gin.HandlerFunc {
	const Bearar = "Bearar "

	return func(ctx *gin.Context) {
		// settings.JwtHAC
		var (
			e      error
			bearar string
			code   string
			err    *errx.ErrX
			data   *ginx.JwtData
		)

		bearar = ctx.GetHeader("Authorization")

		handleError := func() {
			structs.JsonErr(ctx, err)
			ctx.Abort()
		}

		// if bearar[:7] != "Bearar " {
		if !strings.HasPrefix(bearar, Bearar) {
			err = structs.AuthError(fmt.Errorf("...")).WithCode("invalid_token")

			handleError()
			return
		}

		if data, code, e = settings.JwtHMAC.Auth(bearar[len(Bearar):]); e != nil {
			if code == "token_expired" {
				err = structs.AuthError(e).WithMsg("token expired")
			} else {
				err = structs.AuthError(e).WithMsg("invalid token")
			}

			err.WithCode(code)
			handleError()
			return
		}

		ctx.Set("accountId", data.Subject)
		ctx.Set("level", data.Data["level"])
		ctx.Set("platform", data.Data["platform"])
		ctx.Set("tokenId", data.ID)

		for i := range funcs {
			if err = funcs[i](ctx, data); err != nil {
				handleError()
				return
			}
		}

		ctx.Next()
	}
}

package internal

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"backend-api/internal/settings"
	"backend-api/pkg/erri"
	"backend-api/pkg/structs"

	"github.com/d2jvkpn/errx"
	"github.com/d2jvkpn/gotk/ginx"
	"github.com/gin-gonic/gin"
)

type HandleJwt func(context.Context, *ginx.JwtData) *errx.ErrX

func AllowLevels(levels ...string) HandleJwt {
	noAllowed := fmt.Errorf("levels: %v", levels)

	return func(ctx context.Context, data *ginx.JwtData) (err *errx.ErrX) {
		if len(levels) > 0 && !slices.Contains(levels, data.Data["level"]) {
			err = erri.NotPermited(noAllowed).WithCode("no_allowed").WithMsg("no allowed")
			return err
		}

		return nil
	}
}

func AuthCacheUpdateToken(ctx context.Context, d *ginx.JwtData) (err *errx.ErrX) {
	// check if cache token enabled or not internal
	err = settings.CacheUpdateToken(ctx, fmt.Sprintf("%s/%s", d.Subject, d.ID))
	// fmt.Printf("==> AuthCachedToken: %v\n", err)
	if err != nil {
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
			auth   *structs.AuthAccount
		)

		bearar = ctx.GetHeader("Authorization")

		handleError := func() {
			structs.JSONErr(ctx, err)
			ctx.Abort()
		}

		// if bearar[:7] != "Bearar " {
		if !strings.HasPrefix(bearar, Bearar) {
			err = erri.AuthErr(fmt.Errorf("...")).WithCode("invalid_token")

			handleError()
			return
		}

		if data, code, e = settings.JwtHMAC.Auth(bearar[len(Bearar):]); e != nil {
			if code == "token_expired" {
				err = erri.AuthErr(e).WithMsg("token expired")
			} else {
				err = erri.AuthErr(e).WithMsg("invalid token")
			}

			err.WithCode(code)
			handleError()
			return
		}

		auth, e = structs.NewAuthAccount(data.Subject, data.Data["role"], data.ID)
		if e != nil {
			err = erri.AuthErr(e).WithCode("invalid_token").WithMsg("invalid token")

			handleError()
			return
		}
		ctx.Set("Auth", auth)

		structs.GinSetData(ctx, "accountId", auth.AccountId)
		structs.GinSetData(ctx, "level", auth.Level)
		structs.GinSetData(ctx, "tokenId", auth.TokenId)

		for i := range funcs {
			if err = funcs[i](ctx, data); err != nil {
				handleError()
				return
			}
		}

		ctx.Next()
	}
}

package services

import (
	"context"
	"fmt"

	// "net/http"

	"backend-api/internal/biz/biz_user"
	"backend-api/internal/models/mod_user"
	"backend-api/pkg/middlewares"
	"backend-api/pkg/structs"
	"backend-api/pkg/utils"

	"github.com/d2jvkpn/errx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

//	@Summary		Create an account
//	@Description	...
//	@Tags			account::create_account
//	@Accept			json
//	@Produces		json
//	@Param			request								body		mod_user.CreateAccount	true	"account data"
//	@Success		200									{object}	map[string]string
//	@Router			/api/v1/auth/account/create_account	[post]
func createAccount(ctx *gin.Context) {
	var (
		err   *errx.ErrX
		input mod_user.CreateAccount
	)

	if err = BindJSON(ctx, &input); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	if err = input.Do(ctx); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	structs.JsonOK(ctx, gin.H{"accountId": input.Id})
}

//	@Summary		Account login
//	@Description	login in with phone/email and password
//	@Tags			account::login
//	@Accept			json
//	@Produces		json
//	@Param			request						body		biz_user.LoginRequest	true	"login data"
//	@Success		200							{object}	biz_user.LoginResponse
//	@Router			/api/v1/open/account/login	[post]
func accountLogin(ctx *gin.Context) {
	var (
		err    *errx.ErrX
		input  biz_user.LoginRequest
		result *biz_user.LoginResponse

		tracer  trace.Tracer
		span    trace.Span
		otelCtx context.Context
	)

	tracer = otel.Tracer("api.accountLogin")
	otelCtx, span = tracer.Start(ctx, middlewares.GetAPI(ctx))
	// span = trace.SpanFromContext(ctx.Request.Context())
	span.SetAttributes(attribute.String("requestId", ctx.GetString("requestId")))
	defer span.End()
	// fmt.Printf(
	//	"==> TraceID=%s, SpanID=%s\n",
	//	span.SpanContext().TraceID(), span.SpanContext().SpanID(),
	// )

	if err = BindQueryJSON(ctx, &input); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	input.IP, input.TokenId = ctx.ClientIP(), uuid.New().String()

	ctx.Set("platform", input.Platform)
	structs.GinSetData(ctx, "email", input.Email)
	structs.GinSetData(ctx, "phone", input.Phone)
	// ctx and otelCtx share the same key "data"
	structs.GinSetData(ctx, "User-Agent", ctx.GetHeader("User-Agent")) // key=data, pass to biz layer and model layer

	if result, err = biz_user.Login(otelCtx, &input); err != nil {
		structs.JsonErr(ctx, err)
		return
	}
	ctx.Set("tokenId", input.TokenId)
	ctx.Set("accountId", result.Id.String())
	ctx.Set("level", result.Level)

	// ctx.SetCookie("tokenId", tokenId.String(), 300, "/", "localhost", false, true)

	structs.JsonOK(ctx, result)
}

//	@Summary		Account logout
//	@Description	...
//	@Tags			account::logout
//	@Produces		json
//	@Success		200							{object}	structs.ResponseOK
//	@Router			/api/v1/auth/account/logout	[post]
func accountLogout(ctx *gin.Context) {
	var (
		key string
		err *errx.ErrX
	)

	key = fmt.Sprintf("login:%s:%s", ctx.GetString("platform"), ctx.GetString("accountId"))

	if err = biz_user.AccountLogout(ctx, key); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	ctx.Set("skipCacheUpdateToken", true)

	structs.JsonOK(ctx)
}

//	@Summary		Change Password
//	@Description	Change the password of an account
//	@Tags			account::change_password
//	@Accept			json
//	@Produces		json
//	@Param			request									body	mod_user.ChangePassword	true	"oldPasword newPassword"
//	@Router			/api/v1/auth/account/change_password	[post]
func accountChangePassword(ctx *gin.Context) {
	var (
		err       *errx.ErrX
		accountId uuid.UUID
		input     mod_user.ChangePassword

		tracer  trace.Tracer
		span    trace.Span
		otelCtx context.Context
	)

	tracer = otel.Tracer("api.accountChangePassword")
	otelCtx, span = tracer.Start(ctx, middlewares.GetAPI(ctx))
	// span = trace.SpanFromContext(ctx.Request.Context())
	span.SetAttributes(attribute.String("requestId", ctx.GetString("requestId")))
	defer span.End()

	if accountId, err = GetAccountId(ctx); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	if err = BindJSON(ctx, &input); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	err = biz_user.ChangePassword(
		otelCtx, input, accountId,
		fmt.Sprintf("login:%s:%s", ctx.GetString("platform"), accountId),
	)
	if err != nil {
		structs.JsonErr(ctx, err)
		return
	}
	ctx.Set("skipCacheUpdateToken", true)

	structs.JsonOK(ctx)
}

//	@Summary		Query accounts
//	@Description	...
//	@Tags			account::query
//	@Accept			json
//	@Produces		json
//	@Param			query								query		mod_user.QueryAccounts	true	"parameters"
//	@Success		200									{object}	utils.PageResult[mod_user.Account]
//	@Router			/api/v1/auth/account/query_accounts	[get]
func queryAccounts(ctx *gin.Context) {
	var (
		err    *errx.ErrX
		input  mod_user.QueryAccounts
		result *utils.PageResult[mod_user.Account]
	)

	if err = BindQuery(ctx, &input); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	if result, err = input.Do(ctx); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	structs.JsonOK(ctx, result)
}

//	@Summary		Delete accounts
//	@Description	...
//	@Tags			account::query
//	@Accept			json
//	@Produces		json
//	@Param			accountId								query		[]string		string	"uuids of accounts"
//	@Success		200										{object}	map[string]int	'{"count": 32}'
//	@Router			/api/v1/auth/account/delete_accounts	[get]
func deleteAccounts(ctx *gin.Context) {
	var (
		err        *errx.ErrX
		accountIds []uuid.UUID
		count      int64
	)

	if accountIds, err = QueryUUIDs(ctx, "accountId"); err != nil {
		structs.JsonErr(ctx, err)
	}

	if count, err = mod_user.DeleteAccounts(ctx, accountIds); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	structs.JsonOK(ctx, gin.H{"count": count})
}

package services

import (
	"context"
	// "fmt"
	// "net/http"

	"backend-api/internal/biz/biz_user"
	// "backend-api/internal/models"
	"backend-api/pkg/middlewares"
	"backend-api/pkg/structs"

	"github.com/d2jvkpn/errx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// @Summary		Account login
// @Description	login in with phone/email and password
// @Tags			account::login
// @Accept			json
// @Produces		json
// @Param			request						body		biz_user.LoginRequest	true	"login data"
// @Success		200							{object}	biz_user.LoginResponse
// @Router			/api/v1/open/account/login	[post]
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
	// ctx and otelCtx share the same key "data"
	structs.GinSetData(ctx, "User-Agent", ctx.GetHeader("User-Agent")) // pass to biz layer and model layer

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

// @Summary		Account logout
// @Description	...
// @Tags			account::logout
// @Produces		json
// @Success		200							{object}	ResponseOK
// @Router			/api/v1/auth/account/logout	[post]
func accountLogout(ctx *gin.Context) {
	var (
		err  *errx.ErrX
		auth *structs.AuthAccount
	)

	if auth, err = GetAuthAccount(ctx); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	if err = biz_user.AccountLogout(ctx, auth); err != nil {
		structs.JsonErr(ctx, err)
		return
	}

	structs.JsonOK(ctx)
}

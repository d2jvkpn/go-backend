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
		err     *errx.ErrX
		tokenId uuid.UUID
		auth    *structs.AuthAccount
		input   biz_user.LoginRequest
		result  *biz_user.LoginResponse

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

	if err = BindJSON(ctx, &input); err != nil {
		middlewares.JsonErr(ctx, err)
		return
	}

	tokenId = uuid.New()
	input.IP, input.TokenId = ctx.ClientIP(), tokenId.String()

	if result, err = biz_user.Login(otelCtx, &input); err != nil {
		middlewares.JsonErr(ctx, err)
		return
	}

	auth = &structs.AuthAccount{AccountId: result.Id, Level: result.Level, TokenId: tokenId}

	ctx.Set("AuthAccount", auth)
	structs.ContextSetData(ctx, "level", result.Level)
	structs.ContextSetData(ctx, "userAgent", ctx.GetHeader("User-Agent"))

	// ctx.SetCookie("tokenId", tokenId.String(), 300, "/", "localhost", false, true)

	middlewares.JsonOK(ctx, result)
}

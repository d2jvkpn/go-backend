package biz_user

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"backend-api/internal/models/mod_user"
	"backend-api/internal/settings"
	"backend-api/pkg/structs"

	"github.com/d2jvkpn/errx"
	"github.com/google/uuid"
	// "go.opentelemetry.io/otel/attribute"
	"github.com/d2jvkpn/gotk/ginx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type LoginRequest struct {
	mod_user.AccountLogin
	Platform string `json:"-" form:"_platform" extensions:"x-order=04"`

	TokenId string `json:"-" form:"-" swaggerignore:"true"`
	IP      string `json:"-" form:"-" swaggerignore:"true"`
}

type LoginResponse struct {
	*mod_user.Account

	TokenId string `json:"tokenId" extensions:"x-order=21"`
	Token   string `json:"token" extensions:"x-order=22"`
	// unix 时间戳
	ExpiresAt int64 `json:"expiresAt" example:"1715008592" extensions:"x-order=23"`
}

func Login(ctx context.Context, input *LoginRequest) (result *LoginResponse, err *errx.ErrX) {
	var (
		e       error
		token   string
		values  url.Values
		account *mod_user.Account
		data    ginx.JwtData

		tracer  trace.Tracer
		span    trace.Span
		otelCtx context.Context
	)

	// TODO: check input.Platform
	tracer = otel.Tracer("biz_user.Login")

	// 1.
	ctx = context.WithValue(ctx, "platform", input.Platform)
	ctx = context.WithValue(ctx, "tokenId", input.TokenId)

	otelCtx, span = tracer.Start(ctx, "AccountLogin")
	account, err = input.AccountLogin.Do(otelCtx)
	span.End()
	if err != nil {
		return nil, err
	}

	// 2.
	data = ginx.JwtData{
		ID:      input.TokenId,
		Subject: account.Id.String(),
		Data:    map[string]string{"level": account.Level, "platform": input.Platform},
	}

	_, span = tracer.Start(ctx, "JwtHAC.Sign")
	token, e = settings.JwtHMAC.Sign(&data)
	span.End()
	if e != nil {
		return nil, structs.InternalError(e).WithCode("jwt_sign")
	}

	result = &LoginResponse{
		TokenId:   input.TokenId,
		Token:     token,
		ExpiresAt: data.ExpiresAt,
		Account:   account,
	}

	// 3.
	values = make(url.Values, 5)
	values.Add("issuedAt", strconv.FormatInt(data.IssuedAt, 10))
	values.Add("level", account.Level)
	values.Add("ip", input.IP)
	values.Add("tokenId", data.ID)
	values.Add("expiresAt", strconv.FormatInt(data.ExpiresAt, 10))

	_, span = tracer.Start(ctx, "CacheSetToken")
	err = settings.CacheSetToken(
		ctx,
		fmt.Sprintf("login:%s:%s", data.Data["platform"], data.Subject),
		values.Encode(),
	)
	span.End()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func AccountLogout(ctx context.Context, key string) (err *errx.ErrX) {
	err = settings.CacheRemoveToken(ctx, key)

	if err != nil {
		return err
	}

	return nil
}

func ChangePassword(ctx context.Context, input mod_user.ChangePassword, accountId uuid.UUID, key string) (
	err *errx.ErrX) {
	if err = input.Do(ctx, accountId); err != nil {
		return nil
	}

	if err = settings.CacheRemoveToken(ctx, key); err != nil {
		settings.Logger.Named("biz_user").Error("CacheRemoveToken", zap.Any("error", &err))
		err = nil
	}

	return nil
}

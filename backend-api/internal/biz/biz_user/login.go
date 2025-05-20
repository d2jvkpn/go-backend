package biz_user

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"backend-api/internal/models/mod_user"
	"backend-api/internal/settings"
	"backend-api/pkg/erri"
	"backend-api/pkg/structs"

	"github.com/d2jvkpn/errx"
	// "go.opentelemetry.io/otel/attribute"
	"github.com/d2jvkpn/gotk/ginx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type LoginRequest struct {
	mod_user.AccountLogin

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

	tracer = otel.Tracer("biz_user.Login")

	// 1.
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
		Data:    map[string]string{"level": account.Level},
	}

	_, span = tracer.Start(ctx, "JwtHAC.Sign")
	token, e = settings.JwtHMAC.Sign(&data)
	span.End()
	if e != nil {
		return nil, erri.InternalErr(e).WithCode("jwt_sign")
	}

	result = &LoginResponse{
		TokenId:   input.TokenId,
		Token:     token,
		ExpiresAt: data.ExpiresAt,
		Account:   account,
	}

	// 3.
	values = make(url.Values, 3)
	values.Add("issuedAt", strconv.FormatInt(data.IssuedAt, 10))
	values.Add("level", account.Level)
	values.Add("ip", input.IP)

	_, span = tracer.Start(ctx, "CacheSetLogin")
	err = settings.CacheSetLogin(
		ctx,
		fmt.Sprintf("%s/%s", data.Subject, data.ID),
		values.Encode(),
	)
	span.End()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func AccountLogout(ctx context.Context, auth *structs.AuthAccount) (err *errx.ErrX) {
	err = settings.CacheRemoveLogin(
		ctx,
		fmt.Sprintf("%s/%s", auth.AccountId, auth.TokenId),
	)

	if err != nil {
		return err
	}

	return nil
}

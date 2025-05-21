package structs

import (
	"context"
	"fmt"

	"github.com/d2jvkpn/gotk/ginx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthAccount struct {
	AccountId uuid.UUID `json:"accountId"`
	TokenId   uuid.UUID `json:"tokenId"`
	Level     string    `json:"level"`
	Platform  string    `json:"platform"`
}

func NewAuthAccount(id, level, tokenId, platform string) (item *AuthAccount, err error) {
	item = new(AuthAccount)

	if item.AccountId, err = uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	item.Level, item.Platform = level, platform
	if item.TokenId, err = uuid.Parse(tokenId); err != nil {
		return nil, fmt.Errorf("invalid tokenId: %w", err)
	}

	return item, nil
}

func GinSetData(ctx *gin.Context, key string, value any) {
	var (
		e    error
		data map[string]any
	)

	if data, e = ginx.Get[map[string]any](ctx, "data"); e != nil {
		data = make(map[string]any, 1)
		ctx.Set("data", data)
	}

	data[key] = value
}

func ContextSetData(ctx context.Context, key string, value any) context.Context {
	var (
		ok   bool
		data map[string]any
	)

	if data, ok = ctx.Value("data").(map[string]any); !ok {
		data = make(map[string]any, 1)
		ctx = context.WithValue(ctx, "data", data)
	}

	data[key] = value

	return ctx
}

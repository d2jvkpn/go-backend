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
	Level     string    `json:"level"`
	TokenId   uuid.UUID `json:"tokenId"`
}

func NewAuthAccount(id, level, tokenId string) (item *AuthAccount, err error) {
	item = new(AuthAccount)

	if item.AccountId, err = uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	item.Level = level
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

	if data, e = ginx.Get[map[string]any](ctx, "Data"); e != nil {
		data = make(map[string]any, 1)
		ctx.Set("Data", data)
	}

	data[key] = value
}

func ContextSetData(ctx context.Context, key string, value any) context.Context {
	var (
		ok   bool
		data map[string]any
	)

	if data, ok = ctx.Value("Data").(map[string]any); !ok {
		data = make(map[string]any, 1)
		ctx = context.WithValue(ctx, "Data", data)
	}

	data[key] = value

	return ctx
}

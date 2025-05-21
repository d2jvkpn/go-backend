package settings

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"backend-api/pkg/structs"

	"github.com/d2jvkpn/errx"
)

func cacheTokenEnabled() bool {
	// fmt.Printf("==> LoginTokenEnabled: %t\n", Config.GetDuration("jwt.interval") > 0)
	return Config.GetDuration("jwt.interval") > 0
}

func CacheSetToken(ctx context.Context, key, value string) (err *errx.ErrX) {
	if !cacheTokenEnabled() {
		return nil
	}

	var (
		e          error
		expiration time.Duration
	)

	expiration = Config.GetDuration("jwt.interval")
	// fmt.Println("==> LoginTokenSet:", expiration)

	// key := fmt.Sprintf("%s/%s/%s", pkg.CACHE_LoginTokenV1,
	if e = Redis.Set(ctx, key, value, expiration).Err(); e != nil {
		return structs.InternalError(e).WithCode("cache_set_token")
	}

	return nil
}

func CacheUpdateToken(ctx context.Context, key, tokenId string) (err *errx.ErrX) {
	if !cacheTokenEnabled() {
		return nil
	}

	var (
		ok         bool
		e          error
		value      string
		query      url.Values
		expiration time.Duration
	)

	const msg = "Please log in again"
	if value, e = Redis.Get(ctx, key).Result(); e != nil {
		return structs.AuthError(e).WithCode("cache_no_token").WithMsg(msg)
	}

	if query, e = url.ParseQuery(value); e != nil {
		return structs.InternalError(e).WithCode("cache_parse_value")
	}
	if query.Get("tokenId") != tokenId {
		return structs.AuthError(fmt.Errorf("token is expired")).WithCode("cache_token_is_expired").WithMsg(msg)
	}

	expiration = Config.GetDuration("jwt.interval")
	// fmt.Printf("==> CacheLoginTokenUpdate 1: %s, %v\n", key, expiration)

	// problem with concurrency: not ok if expiration == expiration
	// ok, e = Redis.ExpireGT(ctx, key, expiration).Result()
	ok, e = Redis.Expire(ctx, key, expiration).Result()
	// fmt.Printf("==> CacheLoginTokenUpdate 2: %t, %v\n", ok, e)
	if e != nil {
		return structs.InternalError(e).WithCode("cache_update_token")
	}
	if !ok {
		return structs.AuthError(fmt.Errorf("expire token")).WithCode("cache_expire_token").WithMsg(msg)
	}

	return nil
}

func CacheRemoveToken(ctx context.Context, key string) (err *errx.ErrX) {
	if !cacheTokenEnabled() {
		return nil
	}

	var e error
	// fmt.Println("???", key)

	if e = Redis.Del(ctx, key).Err(); e != nil {
		return structs.InternalError(e).WithCode("cache_delete_token")
	}

	return nil
}

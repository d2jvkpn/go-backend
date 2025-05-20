package settings

import (
	"context"
	// "fmt"

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

	return nil
}

func CacheRemoveToken(ctx context.Context, key string) (err *errx.ErrX) {
	if !cacheTokenEnabled() {
		return nil
	}

	return nil
}

func CacheUpdateToken(ctx context.Context, key string) (err *errx.ErrX) {
	if !cacheTokenEnabled() {
		return nil
	}

	return nil
}

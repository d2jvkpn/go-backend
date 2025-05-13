package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedisClient(vp *viper.Viper) (client *redis.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client = redis.NewClient(&redis.Options{
		Addr:                  vp.GetString("addr"),
		Username:              vp.GetString("username"),
		Password:              vp.GetString("password"),
		DB:                    vp.GetInt("db"),
		Protocol:              vp.GetInt("protocol"),
		ContextTimeoutEnabled: true,
		MinIdleConns:          vp.GetInt("min_idle_conns"),
		MaxIdleConns:          vp.GetInt("max_idle_conns"),
		MaxActiveConns:        vp.GetInt("max_active_conns"),
		// ConnMaxIdleTime: 30*time.Minute,
		// ConnMaxLifetime: -1,
	})

	if err = client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis.Ping: %w", err)
	}

	return client, nil
}

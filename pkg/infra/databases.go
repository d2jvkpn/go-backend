package infra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

/*
dsn:
- postgres://{USERANME}:{PASSWORD}@tcp({IP})/{DATABASE}?sslmode=disable
- "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai"
*/
func PgConnect(vp *viper.Viper, release bool) (gormDB *gorm.DB, sqlDB *sql.DB, err error) {
	var conf *gorm.Config

	conf = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}

	if release {
		conf.Logger = logger.Default.LogMode(logger.Silent)
	}

	if gormDB, err = gorm.Open(postgres.Open(vp.GetString("dsn")), conf); err != nil {
		return nil, nil, err
	}

	if !release {
		gormDB = gormDB.Debug()
	}

	// Get generic database object sql.DB to use its functions
	if sqlDB, err = gormDB.DB(); err != nil {
		return nil, nil, err
	}

	if v := vp.GetInt("max_idle_conns"); v > 0 {
		sqlDB.SetMaxIdleConns(v)
	}
	if v := vp.GetInt("max_open_conns"); v > 0 {
		sqlDB.SetMaxOpenConns(v)
	}
	if v := vp.GetDuration("conn_max_idle_time"); v > 0 {
		sqlDB.SetConnMaxLifetime(v)
	}
	if v := vp.GetDuration("conn_max_lifetime"); v > 0 {
		sqlDB.SetConnMaxLifetime(v)
	}

	return gormDB, sqlDB, err
}

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

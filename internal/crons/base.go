package crons

import (
	// "fmt"
	"database/sql"
	"log/slog"

	"github.com/d2jvkpn/gotk"
	"github.com/redis/go-redis/v9"
	// "go.uber.org/zap"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var (
	_SLogger *slog.Logger
	_Logger  *gotk.ZapLogger

	_DB      *sql.DB
	_GORM_PG *gorm.DB
	_Redis   *redis.Client

	_Cron *cron.Cron
)

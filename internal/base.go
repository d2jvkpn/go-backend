package internal

import (
	"context"
	"embed"
	// "fmt"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/trace_error"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	//go:embed static
	_Static embed.FS
	//go:embed templates
	_Templates embed.FS

	_SLogger *slog.Logger
	_Logger  *zap.Logger

	_InternalServer *http.Server
	_HttpServer     *http.Server
	_RPCServer      *RPCServer

	_DB      *sql.DB
	_GORM_PG *gorm.DB
	_Redis   *redis.Client
	// _GORM_MySQL *gorm.DB
	_Tickers    []*gotk.Ticker
	_APIMetrics []func(string, float64, *trace_error.Error)

	_CloseOtelTracing = func(context.Context) error { return nil }
	_CloseOtelMetrics = func(context.Context) error { return nil }
)

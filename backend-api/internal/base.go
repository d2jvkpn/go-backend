package internal

import (
	"context"
	"embed"
	// "fmt"
	"database/sql"
	"log/slog"
	"net/http"

	"backend-api/internal/rpc"

	"github.com/d2jvkpn/gotk"
	es "github.com/elastic/go-elasticsearch/v9"
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

	_DB      *sql.DB
	_GORM_PG *gorm.DB
	_Redis   *redis.Client
	_ES      *es.Client
	// _GORM_MySQL *gorm.DB
	_Tickers   []*gotk.Ticker
	_APIMeters []func(string, float64, []string)

	_InternalServer *http.Server
	_HttpServer     *http.Server
	_RPCServer      *rpc.RPCServer

	_CloseOtelTrace func(context.Context) error = nil // func(context.Context) error { return nil }
	_CloseOtelMeter func(context.Context) error = nil // func(context.Context) error { return nil }
)

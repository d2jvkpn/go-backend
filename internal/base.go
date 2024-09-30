package internal

import (
	"context"
	"embed"
	// "fmt"
	"log/slog"
	"net/http"

	"go.uber.org/zap"
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

	_CloseOtelTracing = func(context.Context) error { return nil }
	_CloseOtelMetrics = func(context.Context) error { return nil }
)

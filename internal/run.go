package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/d2jvkpn/go-backend/internal/settings"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Run(project *viper.Viper) (errch chan error, err error) {
	var (
		httpListener     net.Listener
		internalListener net.Listener
		grpcListener     net.Listener
	)

	defer func() {
		if err == nil {
			return
		}

		if httpListener != nil {
			_SLogger.Info("close http listener")
			err = errors.Join(err, httpListener.Close())
		}

		if internalListener != nil {
			_SLogger.Info("close internal listener")
			err = errors.Join(err, internalListener.Close())
		}

		if grpcListener != nil {
			_SLogger.Info("close grpc listener")
			err = errors.Join(err, grpcListener.Close())
		}
	}()

	_Logger.Info("run", zap.Any("meta", project.GetStringMap("meta")))

	err = gotk.ConcRunErr(
		func() (err error) {
			addr := project.GetString("meta.http_addr")
			_SLogger.Debug("http listen", "address", addr)

			if httpListener, err = net.Listen("tcp", addr); err != nil {
				return fmt.Errorf("http net.Listen: %w", err)
			}
			return nil
		},
		func() (err error) {
			addr := project.GetString("meta.internal_addr")
			_SLogger.Debug("internal listen", "address", addr)

			internalListener, err = net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("internal net.Listen: %w", err)
			}
			return nil
		},
		func() (err error) {
			addr := project.GetString("meta.grpc_addr")
			_SLogger.Debug("grpc listen", "address", addr)

			if grpcListener, err = net.Listen("tcp", addr); err != nil {
				return fmt.Errorf("grpc net.Listen: %w", err)
			}
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	_SLogger.Debug("serve...")
	errch = make(chan error, 3)
	go ServeHTTP(httpListener, errch)
	go ServeInternal(internalListener, errch)
	go ServeGrpc(grpcListener, errch)

	return errch, nil
}

func Exit() (err error) {
	var e error

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	joinErr := func(e error) {
		err = errors.Join(err, e)
	}

	// 1. stop http.server
	shutdownHttp := func() (e error) {
		if _HttpServer == nil {
			return nil
		}

		_SLogger.Debug("shutdown http server")
		if e = _HttpServer.Shutdown(ctx); e != nil {
			_Logger.Error("shutdown http server", zap.String("error", e.Error()))
		}
		return e
	}

	shutdownInternal := func() (e error) {
		if _InternalServer == nil {
			return
		}

		_SLogger.Debug("shutdown internal server")
		if e = _InternalServer.Shutdown(ctx); e != nil {
			_Logger.Error("shutdown internal server", zap.String("error", e.Error()))
		}
		return e
	}

	shutdownGrpc := func() (e error) {
		if _RPCServer == nil {
			return nil
		}
		_SLogger.Debug("shutdown grpc server")
		_RPCServer.Server.GracefulStop()
		return nil
	}

	joinErr(gotk.ConcRunErr(shutdownHttp, shutdownInternal, shutdownGrpc))

	// 2. close otel
	e = gotk.ConcRunErr(
		func() error {
			if _CloseOtelTracing == nil {
				return nil
			}
			_SLogger.Debug("shutdown tracing")
			return _CloseOtelTracing(ctx)
		},
		func() error {
			if _CloseOtelMetrics == nil {
				return nil
			}
			_SLogger.Debug("shutdown metrics")
			return _CloseOtelMetrics(ctx)
		},
		func() error {
			if settings.WsServer == nil {
				return nil
			}

			_SLogger.Debug("shutdown websocket")
			settings.WsServer.Shutdown()
			return nil
		},
	)
	if e != nil {
		_Logger.Error("close otel", zap.String("error", e.Error()))
		joinErr(e)
	}

	// 3. close databases: postgres and redis
	e = gotk.ConcRunErr(
		func() error {
			if _Redis == nil {
				return nil
			}

			_SLogger.Debug("shutdown redis")
			return _Redis.Close()
		},
		func() error {
			if _DB == nil {
				return nil
			}

			_SLogger.Debug("shutdown postgres")
			return _DB.Close()
		},
	)
	if e != nil {
		_Logger.Error("close databases", zap.String("error", e.Error()))
		joinErr(e)
	}

	// 4. close logger
	if settings.Logger != nil {
		_SLogger.Debug("shutdown logger")

		if err == nil {
			_Logger.Info("exit")
		} else {
			_Logger.Error("exit", zap.String("error", e.Error()))
		}

		if e = settings.Logger.Down(); e != nil {
			joinErr(e)
		}
	}

	return err
}

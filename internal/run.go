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

func Load(project *viper.Viper) (err error) {
	var (
		appName string
		release bool
		config  *viper.Viper
	)

	release = project.GetBool("meta.release")
	appName = project.GetString("app_name")

	config, err = gotk.LoadYamlConfig(project.GetString("meta.config"), "config")
	if err != nil {
		return err
	}
	config.SetDefault("prometheus", map[string]any{})
	config.SetDefault("opentelemetry", map[string]any{})

	// 1. Log
	if err = SetupLog(release, project.GetString("app_name")); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			Exit()
		}
	}()

	// 2. databases: postgres, redis
	err = gotk.ConcRunErr(
		func() (err error) {
			_GORM_PG, _DB, err = PgConnect(config.Sub("postgres"), release)
			return err
		},
		func() (err error) {
			_Redis, err = NewRedisClient(config.Sub("redis"))
			return err
		},
	)
	if err != nil {
		return err
	}

	// 3. cloud
	err = gotk.ConcRunErr(
		func() (err error) {
			return SetupOtelMetrics(appName, config)
		},
		func() (err error) {
			return SetupOtelTrace(appName, config)
		},
	)
	if err != nil {
		return err
	}

	// 4. http server
	if err = SetupHttp(release, config); err != nil {
		return err
	}

	// 5. internal server
	if err = SetupInternal(config, project.GetStringMap("meta")); err != nil {
		return err
	}

	// 6. grpc server
	if err = SetupGrpc(config); err != nil {
		return err
	}

	return nil
}

func Run(project *viper.Viper) (errch chan error, err error) {
	var (
		httpListener     net.Listener
		grpcListener     net.Listener
		internalListener net.Listener
	)

	defer func() {
		if err == nil {
			return
		}

		if httpListener != nil {
			err = errors.Join(err, httpListener.Close())
		}

		if internalListener != nil {
			err = errors.Join(err, internalListener.Close())
		}

		if grpcListener != nil {
			err = errors.Join(err, grpcListener.Close())
		}
	}()

	_Logger.Info("run", zap.Any("meta", project.GetStringMap("meta")))
	err = gotk.ConcRunErr(
		func() (err error) {
			httpListener, err = net.Listen("tcp", project.GetString("meta.http_addr"))
			if err != nil {
				return fmt.Errorf("http net.Listen: %w", err)
			}
			return nil
		},
		func() (err error) {
			internalListener, err = net.Listen("tcp", project.GetString("meta.internal_addr"))
			if err != nil {
				return fmt.Errorf("internal net.Listen: %w", err)
			}
			return nil
		},
		func() (err error) {
			grpcListener, err = net.Listen("tcp", project.GetString("meta.grpc_addr"))
			if err != nil {
				return fmt.Errorf("grpc net.Listen: %w", err)
			}
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

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
	if _HttpServer != nil {
		_SLogger.Warn("shutdown http server")
		if e = _HttpServer.Shutdown(ctx); e != nil {
			_Logger.Error("shutdown http server", zap.String("error", e.Error()))
			joinErr(e)
		}
	}

	// 2. stop internal server
	if _InternalServer != nil {
		_SLogger.Warn("shutdown internal server")
		if e = _InternalServer.Shutdown(ctx); e != nil {
			_Logger.Error("shutdown internal server", zap.String("error", e.Error()))
			joinErr(e)
		}
	}

	// 3. stop grpc sever
	if _RPCServer != nil {
		_SLogger.Warn("shutdown grpc server")
		_RPCServer.Server.GracefulStop()
	}

	// 4. close otel
	e = gotk.ConcRunErr(
		func() error { return _CloseOtelTracing(ctx) },
		func() error { return _CloseOtelMetrics(ctx) },
	)
	if e != nil {
		_Logger.Error("close otel", zap.String("error", e.Error()))
		joinErr(e)
	}

	// 5. close databases: postgres and redis
	e = gotk.ConcRunErr(
		func() error {
			if _Redis == nil {
				return nil
			}
			return _Redis.Close()
		},
		func() error {
			if _DB == nil {
				return nil
			}
			return _DB.Close()
		},
	)
	if e != nil {
		_Logger.Error("close databases", zap.String("error", e.Error()))
		joinErr(e)
	}

	// 6. close logger
	if settings.Logger != nil {
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

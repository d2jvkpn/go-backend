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
		release bool
		config  *viper.Viper
	)

	config, err = gotk.LoadYamlConfig(project.GetString("meta.config"), "config")
	if err != nil {
		return err
	}
	release = project.GetBool("meta.release")

	// 1. Log
	if err = SetupLog(release, project.GetString("app_name")); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			Shutdown()
		}
	}()

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// 4. http server
	if err = SetupHttp(release, config); err != nil {
		return err
	}

	// 5. internal server
	SetupInternal(config, project.GetStringMap("meta"))

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

		if grpcListener != nil {
			err = errors.Join(err, grpcListener.Close())
		}

		if internalListener != nil {
			err = errors.Join(err, internalListener.Close())
		}
	}()

	httpListener, err = net.Listen("tcp", project.GetString("meta.http_addr"))
	if err != nil {
		return nil, fmt.Errorf("http net.Listen: %w", err)
	}

	grpcListener, err = net.Listen("tcp", project.GetString("meta.grpc_addr"))
	if err != nil {
		return nil, fmt.Errorf("grpc net.Listen: %w", err)
	}

	internalListener, err = net.Listen("tcp", project.GetString("meta.internal_addr"))
	if err != nil {
		return nil, fmt.Errorf("internal net.Listen: %w", err)
	}

	errch = make(chan error, 3)
	go ServeHTTP(httpListener, errch)
	go ServeInternal(internalListener, errch)
	go ServeGrpc(grpcListener, errch)

	return errch, nil
}

func Shutdown() (err error) {
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

	// 5. databases

	// 6. close logger
	if settings.Logger != nil {
		if e = settings.Logger.Down(); e != nil {
			_Logger.Error("shutdown logger", zap.String("error", e.Error()))
			joinErr(e)
		}
	}

	// 7. end
	if err == nil {
		_Logger.Warn("end")
	} else {
		_Logger.Error("end", zap.Any("error", &err))
	}

	return err
}

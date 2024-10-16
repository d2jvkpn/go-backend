package internal

import (
	// "context"
	// "fmt"

	"github.com/d2jvkpn/go-backend/internal/rpc"
	"github.com/d2jvkpn/go-backend/internal/settings"
	"github.com/d2jvkpn/go-backend/internal/ws"
	"github.com/d2jvkpn/go-backend/pkg/infra"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/cloud"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/metric"
)

func Load(project *viper.Viper) (err error) {
	var (
		appName string
		release bool
		config  *viper.Viper
	)

	appName = project.GetString("app_name") + ".api"
	release = project.GetBool("meta.release")

	config, err = gotk.LoadYamlConfig(project.GetString("meta.config"), "config")
	if err != nil {
		return err
	}

	config.SetDefault("prometheus", map[string]any{})
	config.SetDefault("opentelemetry", map[string]any{})

	grpcConfig := config.Sub("grpc")
	grpcConfig.Set("trace", config.GetBool("opentelemetry.trace"))
	grpcConfig.Set("meter", config.GetBool("opentelemetry.meter"))

	// 1. Log
	if err = SetupLog(appName, release); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			Exit()
		}
	}()

	otelConfig := config.Sub("opentelemetry")

	// 2. databases(postgres, redis) and otel(tracer and meter)
	err = gotk.ConcRunErr(
		func() (err error) {
			_SLogger.Debug("connect to postgres")
			_GORM_PG, _DB, err = infra.PgConnect(config.Sub("postgres"), release)
			return err
		},
		func() (err error) {
			_SLogger.Debug("connect to redis")
			_Redis, err = infra.NewRedisClient(config.Sub("redis"))
			return err
		},
		func() (err error) {
			if !otelConfig.GetBool("trace") {
				return nil
			}

			_SLogger.Debug("setup otel trace")
			_CloseOtelTrace, err = cloud.OtelTraceGrpc(appName, otelConfig)
			return err
		},
		func() (err error) {
			if !otelConfig.GetBool("meter") {
				return nil
			}

			_SLogger.Debug("setup otel meter")
			_CloseOtelMeter, err = cloud.OtelMeterGrpc(appName, otelConfig, false)

			return err
		},
	)
	if err != nil {
		return err
	}

	// 4. metrcs
	if otelConfig.GetBool("meter") {
		var (
			meter     otelmetric.Meter
			otelMeter func(string, float64, []string)
		)

		meter = otel.GetMeterProvider().Meter(appName)

		// println("==> SetupDBStatsOtel")
		if err = cloud.SetupDBStatsOtel(_DB, meter); err != nil {
			return err
		}

		otelMeter, err = cloud.OtelMeterHttp(meter, []string{"code", "kind"})
		if err != nil {
			return err
		}

		_APIMeters = append(_APIMeters, otelMeter)
	}

	// 5. servers
	// http server
	_SLogger.Debug("setup http")
	if err = SetupHttp(release, config); err != nil {
		return err
	}

	// internal server
	_SLogger.Debug("setup internal")
	if err = SetupInternal(config, project.GetStringMap("meta")); err != nil {
		return err
	}

	// grpc server
	_SLogger.Debug("setup grpc")
	if _RPCServer, err = rpc.NewRPCServer(config); err != nil {
		return err
	}

	settings.WsServer = ws.NewServer(settings.Logger.Named("websocket"))

	return nil
}

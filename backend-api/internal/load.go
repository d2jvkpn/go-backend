package internal

import (
	"context"
	// "fmt"
	"time"

	"backend-api/internal/models/mod_user"
	"backend-api/internal/rpc"
	"backend-api/internal/settings"
	"backend-api/internal/ws"
	"backend-api/pkg/infra"
	"backend-api/pkg/utils"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/cloud"
	"github.com/d2jvkpn/gotk/ginx"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/metric"
)

func Load(project *viper.Viper) (err error) {
	var (
		appName string
		release bool
		ctx     context.Context
		cancel  func()
		config  *viper.Viper
	)

	settings.Project = project
	appName = project.GetString("app_name")
	release = project.GetBool("meta.release")
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	config, err = gotk.LoadYamlConfig(project.GetString("meta.config"), "config")
	if err != nil {
		return err
	}

	config.SetDefault("init", map[string]any{})
	config.SetDefault("prometheus", map[string]any{})
	config.SetDefault("opentelemetry", map[string]any{})
	settings.Config = config

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

	if settings.Captcha, err = utils.CaptchaFromViper(config.Sub("captcha")); err != nil {
		return err
	}

	settings.JwtHMAC, err = ginx.NewJwtHMAC(config.Sub("jwt"), appName)
	if err != nil {
		return err
	}

	otelConfig := config.Sub("opentelemetry")

	// 2. databases(postgres & redis) and otel(tracer & meter)
	err = gotk.ConcRunErr(
		func() (err error) {
			_SLogger.Debug("connecting to postgres")
			_GORM_PG, _DB, err = infra.PgConnect(config.Sub("postgres"), release)

			return err
		},
		func() (err error) {
			_SLogger.Debug("connecting to redis")
			_Redis, err = infra.NewRedisClient(config.Sub("redis"))
			settings.Redis = _Redis
			return err
		},
		func() (err error) {
			_SLogger.Debug("connecting to elasticsearch")
			_ES, err = infra.NewEsClient(config.Sub("elasticsearch"))
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

	// 4. Initialize mod_user
	if err = mod_user.Init(ctx, _GORM_PG); err != nil {
		return err
	}
	if err = mod_user.InitializeDefaultAccount(ctx, config.Sub("init")); err != nil {
		return err
	}

	// 5. metrics
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

		otelMeter, err = cloud.OtelMeterHttp(meter, []string{"kind", "code"})
		if err != nil {
			return err
		}

		_APIMeters = append(_APIMeters, otelMeter)
	}

	// 6. servers
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

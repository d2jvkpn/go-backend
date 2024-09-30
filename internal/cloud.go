package internal

import (
	// "fmt"

	"github.com/d2jvkpn/gotk/cloud"
	"github.com/d2jvkpn/gotk/trace_error"
	"github.com/spf13/viper"
	otelmetric "go.opentelemetry.io/otel/metric"
)

func SetupOtelTrace(appName string, config *viper.Viper) (err error) {
	var otelConfig *viper.Viper

	otelConfig = config.Sub("opentelemetry")

	if !otelConfig.GetBool("trace") {
		return nil
	}

	if _CloseOtelTracing, err = cloud.OtelTracingGrpc(appName, otelConfig); err != nil {
		return err
	}

	return nil
}

func SetupOtelMetrics(appName string, config *viper.Viper) (err error) {
	var (
		otelConfig  *viper.Viper
		meter       otelmetric.Meter
		otelMetrics func(string, float64, *trace_error.Error)
	)

	otelConfig = config.Sub("opentelemetry")

	if !otelConfig.GetBool("metrics") {
		return nil
	}

	if _DB == nil {
		return nil
	}

	meter, _CloseOtelMetrics, err = cloud.OtelMetricsGrpc(appName, otelConfig, false)
	if err != nil {
		return nil
	}

	// println("==> SetupDBStatsOtel")
	if err = cloud.SetupDBStatsOtel(_DB, meter); err != nil {
		return err
	}

	if otelMetrics, err = cloud.OtelMetricsAPI(meter); err != nil {
		return err
	}

	_APIMetrics = append(_APIMetrics, otelMetrics)

	return nil
}

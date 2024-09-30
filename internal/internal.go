package internal

import (
	"errors"
	"net"
	// "fmt"
	"net/http"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/ginx"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	// "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupInternal(config *viper.Viper, meta map[string]any) (err error) {
	var (
		promConfig *viper.Viper
		pprofs     map[string]http.HandlerFunc
		engine     *gin.Engine
		router     *gin.RouterGroup
	)

	promConfig = config.Sub("prometheus")

	engine = gin.New()
	engine.Use(gin.Recovery())
	// engine = gin.Default()
	engine.RedirectTrailingSlash = true

	// engine.NoRoute(...) // TODO

	router = &engine.RouterGroup

	router.GET("/healthz", ginx.Healthz)
	router.GET("/meta", ginx.JSONStatic(meta))

	if promConfig != nil && promConfig.GetBool("enabled") { // !promConfig.GetBool("external")
		router.GET(
			promConfig.GetString("path"),
			// gin.WrapH(promhttp.Handler()),
			gin.WrapH(promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})),
		)
	}
	// mux := http.NewServeMux()
	// mux.Handle(p, promhttp.Handler())

	pprofs = gotk.PprofHandlerFuncs()
	for _, k := range gotk.PprofFuncKeys() {
		router.GET("/pprof/"+k, gin.WrapH(pprofs[k]))
	}

	_InternalServer = &http.Server{
		Handler: engine, // mux
	}

	return nil
}

func ServeInternal(listener net.Listener, errch chan<- error) {
	_SLogger.Info("internal server is up")

	var e error

	e = _InternalServer.Serve(listener)
	_InternalServer = nil // tag as closed

	if e != nil && !errors.Is(e, http.ErrServerClosed) { // e != http.ErrServerClosed
		_Logger.Error("internal server has been shutdown", zap.String("error", e.Error()))
		errch <- e
	} else {
		_Logger.Warn("internal server has been shutdown")
		errch <- nil
	}

	return
}

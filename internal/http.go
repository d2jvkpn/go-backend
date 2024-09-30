package internal

import (
	"crypto/tls"
	"errors"
	"net"
	// "fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/ginx"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	// "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupHttp(release bool, config *viper.Viper) (err error) {
	var (
		fsys       fs.FS
		templ      *template.Template
		httpConfig *viper.Viper
		cert       tls.Certificate

		router *gin.RouterGroup
		engine *gin.Engine
	)

	httpConfig = config.Sub("http")

	// 1. server
	_HttpServer = &http.Server{
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    2 << 11, // 4K
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      10 * time.Second,
		// Handler:           engine,
		// Addr:              addr,
	}

	if httpConfig.GetBool("tls") {
		certFile, keyFile := httpConfig.GetString("cer"), httpConfig.GetString("key")

		if cert, err = tls.LoadX509KeyPair(certFile, keyFile); err != nil {
			return err
		}

		_HttpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	// 2. engine
	if release {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		// engi.Use(gin.Recovery())
	} else {
		engine = gin.Default()
	}
	engine.RedirectTrailingSlash = false
	// TODO: max body size
	// engine.MaxMultipartMemory = HTTP_MaxMultipartMemory

	// engine.Use(Cors(config.GetString("cors")))
	engine.Use(Cors(httpConfig.GetStringSlice("allow_origins")))

	router = &engine.RouterGroup
	if p := httpConfig.GetString("base_path"); p != "" {
		*router = *(router.Group(p))
	}

	// engine.LoadHTMLGlob("templates/*.templ"), "templates/*/*.html"
	templ, err = template.ParseFS(_Templates, "templates/*.html")
	if err != nil {
		return err
	}
	engine.SetHTMLTemplate(templ)

	// 4. middlwares
	// engine.NoRoute(...) // TODO
	// apiLog = ... // TODO

	// 5. apis and router
	router.GET("/healthz", ginx.Healthz)

	if fsys, err = fs.Sub(_Static, "static"); err != nil {
		return err
	}
	static := router.Group("/static", ginx.CacheControl(60))
	static.StaticFS("/", http.FS(fsys))

	ginx.ServeStaticDir("/site", "./site", false)(router)

	// 4. load api
	// TODO: apiLog

	_HttpServer.Handler = engine

	return nil
}

func SetupInternal(config *viper.Viper, meta map[string]any) {
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

	return
}

func ServeHTTP(listener net.Listener, errch chan<- error) {
	_SLogger.Info("http server is up")

	var e error

	e = _HttpServer.Serve(listener)
	_HttpServer = nil // tag as closed

	if e != nil && !errors.Is(e, http.ErrServerClosed) { // e != http.ErrServerClosed
		_Logger.Error("http server has been shutdown", zap.String("error", e.Error()))
		errch <- e
	} else {
		_Logger.Warn("http server has been shutdown")
		errch <- nil
	}
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

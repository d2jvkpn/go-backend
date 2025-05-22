package internal

import (
	"crypto/tls"
	"errors"
	"net"
	// "fmt"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"backend-api/internal/services"
	"backend-api/internal/settings"
	"backend-api/pkg/middlewares"

	"github.com/d2jvkpn/errx"
	"github.com/d2jvkpn/gotk/ginx"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	// "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupHttp(release bool, config *viper.Viper) (err error) {
	var (
		fsys       fs.FS
		templ      *template.Template
		httpConfig *viper.Viper
		cert       tls.Certificate

		apiLog gin.HandlerFunc
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
		// engi.Use(gin.Recovery()) // custom in middleware
	} else {
		engine = gin.Default()
	}
	engine.RedirectTrailingSlash = false
	// engine.MaxMultipartMemory = HTTP_MaxMultipartMemory // ??

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
	notRoute, _ := json.Marshal(gin.H{"code": "no_route", "kind": "no_route", "msg": "..."})
	engine.NoRoute(func(ctx *gin.Context) {
		time.Sleep(1000 * time.Millisecond)

		ctx.Header("Content-Type", "application/json")
		ctx.Writer.WriteHeader(http.StatusNotFound)
		ctx.Writer.Write(notRoute)
	})

	apiLog = middlewares.NewAPILog(
		settings.Logger.Named("api_log"),
		settings.Logger.Level() == zapcore.DebugLevel,
		func(ctx *gin.Context) ([]string, any) {
			if err, _ := ginx.Get[*errx.ErrX](ctx, "error"); err != nil {
				return []string{err.Code, err.Kind}, err
			}

			return []string{}, nil
		},
		_APIMeters...,
	)

	// 5. apis and router
	router.GET("/healthz", ginx.Healthz)

	if fsys, err = fs.Sub(_Static, "static"); err != nil {
		return err
	}
	static := router.Group("/static", ginx.CacheControl(60))
	static.StaticFS("/", http.FS(fsys))

	ginx.ServeStaticDir("/site", "./data/site", false)(router)

	// 6. load api
	services.LoadOpen(router, apiLog)
	services.LoadAuth(router, apiLog, Auth(AllowLevels() /*, CacheUpdateToken*/))
	// services.LoadWebsocket(router)

	_HttpServer.Handler = engine

	return nil
}

func ServeHTTP(listener net.Listener, errch chan<- error) {
	_SLogger.Debug("http server is up")

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

func Cors(origins []string, maxAges ...time.Duration) gin.HandlerFunc {
	maxAge := 12 * time.Hour
	if len(maxAges) > 0 {
		maxAge = maxAges[0]
	}

	return cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "OPTIONS", "HEAD", "PUT"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Authorization",
			"Cache-Control",
			//"X-CSRF-Token",

			"x-client",
		},
		ExposeHeaders: []string{
			"Content-Type",
			"Content-Length",

			"x-api",
		},
		AllowCredentials: true,
		// AllowOriginFunc:  func(origin string) bool { return origin == "https://github.com" },
		MaxAge: maxAge,
	})
}

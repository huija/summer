package summer

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/huija/summer/conf"
	"github.com/huija/summer/container/pipeline"
	"github.com/huija/summer/logs"
	"github.com/huija/summer/srv"
	"github.com/huija/summer/utils"
	"net/http"
	"net/http/pprof"
	"time"
)

// GinSrv gin srv
var GinSrv *gin.Engine

const (
	GinRun   = Gin + Split + "run"
	GinTrace = Gin + Split + "trace"
	GinPprof = Gin + Split + "pprof"
	GinCors  = Gin + Split + "cors"
)

func ginEngine() (err error) {
	if conf.Config.Srv == nil {
		return errors.New("conf.Srv is nil")
	}

	if conf.Config.Srv.FrameWork != srv.GIN {
		return nil
	}

	if conf.Config.Srv.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	srv.Handler = gin.New()

	// logger & recover
	GinSrv = srv.Handler.(*gin.Engine)
	if logs.Writer == nil {
		return errors.New("logs.Writer is nil")
	}
	GinSrv.Use(gin.RecoveryWithWriter(logs.Writer))

	// static template
	if conf.Config.Srv.HtmlGlob != "" {
		GinSrv.LoadHTMLGlob(conf.Config.Srv.HtmlGlob)
	}
	if conf.Config.Srv.Static != "" {
		GinSrv.Static("/static", conf.Config.Srv.Static)
	}

	logs.SugaredLogger.Debug("gin engine init successfully")
	return ginSetup()
}

func ginSetup() (err error) {
	if AddStage(GinRun, ginRun, pipeline.Runner) == nil {
		return errors.New("pipe: add " + GinRun + " stage failed")
	}

	if conf.Config.Srv.Trace {
		if AddStage(GinTrace, ginTrace, pipeline.Debugger) == nil {
			return errors.New("pipe: add " + GinTrace + " stage failed")
		}
	}

	if conf.Config.Srv.Cors {
		if AddStage(GinCors, ginCors, pipeline.Enhancer) == nil {
			return errors.New("pipe: add " + GinCors + " stage failed")
		}
	}

	if conf.Config.Srv.Pprof {
		if AddStage(GinPprof, ginPProf, pipeline.Enhancer) == nil {
			return errors.New("pipe: add " + GinPprof + " stage failed")
		}
	}

	logs.SugaredLogger.Infof("%+v", GetStage(Gin))
	return nil
}

func ginRun() error {
	if GinSrv == nil {
		return errors.New("pipe: gin srv is null")
	}

	address := conf.Config.Srv.Listen + ":" + conf.Config.Srv.Port
	srv.Server = http.Server{Addr: address, Handler: GinSrv}

	logs.SugaredLogger.Debug("starting server on ", address)
	// TODO gin.Run can not be shutdown https://github.com/gin-gonic/gin/issues/2304
	//return GinSrv.Run(address)
	return srv.Server.ListenAndServe()
}

func ginTrace() error {
	GinSrv.Use(trace(logs.SugaredLogger))

	logs.SugaredLogger.Debug("gin trace init successfully!")
	return nil
}

const TraceID = "traceId"

func trace(logger logs.Logger) gin.HandlerFunc {
	srv.GetTraceId = func(c context.Context, wrap ...string) string {
		ctx := c.(*gin.Context)
		switch len(wrap) {
		case 1:
			return wrap[0] + ctx.GetString(TraceID) + wrap[0]
		case 2:
			return wrap[0] + ctx.GetString(TraceID) + wrap[1]
		default:
			return ctx.GetString(TraceID)
		}
	}

	var createTraceID = func(name string, ip string, timestamp time.Time) string {
		return fmt.Sprintf("%s-%s-%d-%s",
			name, ip, timestamp.Unix(),
			utils.NewObjectIDFromTimestamp(timestamp).Hex())
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Set(TraceID, createTraceID(
			conf.Config.Srv.Name+":"+conf.Config.Srv.Tag,
			conf.Config.Srv.IP, start))

		c.Next()

		logger.Debug(srv.GetTraceId(c, "[", "]"),
			"|", c.Writer.Status(),
			"|", time.Since(start),
			"|", c.ClientIP(),
			"|", c.Request.Method,
			"|", c.Request.URL.Path,
			"|", c.Request.URL.RawQuery,
			"|", c.Request.UserAgent(),
			"|", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}

func ginCors() (err error) {
	var ginCors = func() gin.HandlerFunc {
		c := cors.DefaultConfig()
		c.AllowAllOrigins = true
		c.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "AccessToken", "X-CSRF-Token", "Authorization", "Token"}
		c.AllowCredentials = true
		c.ExposeHeaders = []string{"Access-Control-Allow-Origin"}

		return cors.New(c)
	}

	GinSrv.Use(ginCors())

	logs.SugaredLogger.Debug("gin cors init successfully!")
	return nil
}

func ginPProf() (err error) {
	pp := GinSrv.Group("/debug/pprof")

	var pprofHandler = func(h http.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			h.ServeHTTP(c.Writer, c.Request)
		}
	}

	pp.GET("/", pprofHandler(pprof.Index))
	pp.GET("/cmdline", pprofHandler(pprof.Cmdline))
	pp.GET("/profile", pprofHandler(pprof.Profile))
	pp.POST("/symbol", pprofHandler(pprof.Symbol))
	pp.GET("/symbol", pprofHandler(pprof.Symbol))
	pp.GET("/trace", pprofHandler(pprof.Trace))
	pp.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
	pp.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
	pp.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
	pp.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
	pp.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
	pp.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))

	logs.SugaredLogger.Debug("gin pprof router init successfully!")
	return
}
